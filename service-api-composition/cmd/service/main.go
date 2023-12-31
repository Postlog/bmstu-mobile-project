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

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"

	imageStorageClient "github.com/postlog/mobile-project/service-api-composition/internal/clients/image_storage"
	"github.com/postlog/mobile-project/service-api-composition/internal/config/service"
	createScaleTaskHandler "github.com/postlog/mobile-project/service-api-composition/internal/handlers/create_scale_task"
	getImageHandler "github.com/postlog/mobile-project/service-api-composition/internal/handlers/get_image"
	getScaleResultHandler "github.com/postlog/mobile-project/service-api-composition/internal/handlers/get_scale_result"
	infoHandler "github.com/postlog/mobile-project/service-api-composition/internal/handlers/info"
	saveImageHandler "github.com/postlog/mobile-project/service-api-composition/internal/handlers/save_image"
	scaleResultRepository "github.com/postlog/mobile-project/service-api-composition/internal/repository/scale_result"
	scaleTaskRepository "github.com/postlog/mobile-project/service-api-composition/internal/repository/scale_task"
)

func main() {
	os.Exit(run())
}

func run() (exitCode int) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With("origin", "service")
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err := fmt.Errorf("%v", panicErr)

			logger.ErrorContext(ctx, "unhandled panic in service/main", "panic", err)

			exitCode = 1
		}
	}()

	cfg, err := service.Load()
	if err != nil {
		logger.ErrorContext(ctx, "error loading config", "error", err)
		return 1
	}

	logger.InfoContext(ctx, "app config", "config", cfg)

	imageStorageHTTPClient := http.Client{
		Transport: nil,
		Timeout:   cfg.DependenciesConfig.ServiceImageStorageTimeout,
	}

	imageStorageClientInstance := imageStorageClient.New(cfg.DependenciesConfig.ServiceImageStorageURL, imageStorageHTTPClient)
	if err != nil {
		logger.ErrorContext(ctx, "error initializing image-storage client", "error", err)
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

	createScaleTaskHandlerInstance := createScaleTaskHandler.New(logger, scaleTaskRepo)
	getScaleResultHandlerInstance := getScaleResultHandler.New(logger, scaleResultRepo)
	saveImageHandlerInstance := saveImageHandler.New(logger, imageStorageClientInstance)
	getImageHandlerInstance := getImageHandler.New(logger, imageStorageClientInstance)
	infoHandlerInstance := infoHandler.New(logger)

	router := mux.NewRouter()
	router.Handle("/info", infoHandlerInstance).Methods(http.MethodGet)
	router.Handle("/image/{imageId}", getImageHandlerInstance).Methods(http.MethodGet)
	router.Handle("/image", saveImageHandlerInstance).Methods(http.MethodPost)
	router.Handle("/task/scale/{taskId}", getScaleResultHandlerInstance).Methods(http.MethodGet)
	router.Handle("/task/scale", createScaleTaskHandlerInstance).Methods(http.MethodPost)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerConfig.Port),
		Handler: router,
	}

	go func() {
		signalCh := make(chan os.Signal)
		signal.Notify(signalCh, os.Interrupt)

		select {
		case <-signalCh:
			logger.InfoContext(ctx, "stopping server gracefully")
			_ = server.Shutdown(ctx)
		case <-ctx.Done():
		}
	}()

	logger.InfoContext(ctx, fmt.Sprintf("server starting on port %d", cfg.ServerConfig.Port))
	if err = server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logger.ErrorContext(ctx, "server stopped with error", "error", err)
		}
	}

	logger.InfoContext(ctx, "server stopped")

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
