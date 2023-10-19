package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"

	imageScalerClient "github.com/postlog/mobile-project/service-api-composition/internal/clients/image_scaler"
	"github.com/postlog/mobile-project/service-api-composition/internal/config/workers/image_scaler"
	scaleResultRepository "github.com/postlog/mobile-project/service-api-composition/internal/repository/scale_result"
	scaleTaskRepository "github.com/postlog/mobile-project/service-api-composition/internal/repository/scale_task"
	imageScalerService "github.com/postlog/mobile-project/service-api-composition/internal/service/image_scaler"
)

func main() {
	os.Exit(run())
}

func run() (exitCode int) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With("origin", "image_scaler_worker")
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err := fmt.Errorf("%v", panicErr)

			logger.ErrorContext(ctx, "unhandled panic in workers/image_scaler/main", "panic", err)

			exitCode = 1
		}
	}()

	cfg, err := service.Load()
	if err != nil {
		logger.ErrorContext(ctx, "error loading config", "error", err)
		return 1
	}

	logger.InfoContext(ctx, "app config", "config", cfg)

	imageScalerHTTPClient := http.Client{
		Transport: nil,
		Timeout:   cfg.DependenciesConfig.ServiceImageScalerTimout,
	}

	imageScalerClientInstance := imageScalerClient.New(cfg.DependenciesConfig.ServiceImageScalerURL, imageScalerHTTPClient)
	if err != nil {
		logger.ErrorContext(ctx, "error initializing image-sclaer client", "error", err)
		return 1
	}

	rabbitMQConn, err := amqp.Dial(cfg.RabbitConfig.DSN)
	if err != nil {
		logger.ErrorContext(ctx, "error initializing rabbitMQ connection", "error", err)
		return 1
	}
	defer func() { _ = rabbitMQConn.Close() }()

	db, err := sql.Open("postgres", getPostgresDSN(cfg.PostgresConfig))
	if err != nil {
		logger.ErrorContext(ctx, "error initializing postgres connection", "error", err)
		return 1
	}
	defer func() { _ = db.Close() }()

	scaleTaskRepo := scaleTaskRepository.New(rabbitMQConn)
	scaleResultRepo := scaleResultRepository.New(db)

	imageScalerServiceInstance := imageScalerService.New(logger, scaleResultRepo, scaleTaskRepo, imageScalerClientInstance)

	go func() {
		signalCh := make(chan os.Signal)
		signal.Notify(signalCh, os.Interrupt)

		select {
		case <-signalCh:
			logger.InfoContext(ctx, "stopping worker gracefully")
			ctx.Done()
		case <-ctx.Done():
		}
	}()

	logger.InfoContext(ctx, "image scaler worker starting")
	if err = imageScalerServiceInstance.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		logger.ErrorContext(ctx, "error during running image-scaler service", "error", err)
	}

	logger.InfoContext(ctx, "worker stopped")

	return 0
}

func getPostgresDSN(postgresConfig service.PostgresConfig) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		postgresConfig.Host,
		postgresConfig.Port,
		postgresConfig.User,
		postgresConfig.Password,
		postgresConfig.DBName,
	)
}
