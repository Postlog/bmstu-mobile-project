package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	imageStorageClient "github.com/postlog/mobile-project/service-api-composition/internal/clients/image_storage"
	"github.com/postlog/mobile-project/service-api-composition/internal/config"
	getImageHandler "github.com/postlog/mobile-project/service-api-composition/internal/handlers/get"
	saveImageHandler "github.com/postlog/mobile-project/service-api-composition/internal/handlers/save"
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

	cfg, err := config.Load()
	if err != nil {
		logger.ErrorContext(ctx, "error loading config", "error", err)
		return 1
	}

	logger.InfoContext(ctx, "app config", "config", cfg)

	dependenciesHTTPClient := http.Client{
		Transport: nil,
		Timeout:   cfg.DependenciesConfig.ServiceImageStorageTimeout,
	}

	imageStorageClientInstance := imageStorageClient.New(cfg.DependenciesConfig.ServiceImageStorageURL, dependenciesHTTPClient)
	if err != nil {
		logger.ErrorContext(ctx, "error initializing image-storage client", "error", err)
		return 1
	}

	saveImageHandlerInstance := saveImageHandler.New(logger, imageStorageClientInstance)
	getImageHandlerInstance := getImageHandler.New(logger, imageStorageClientInstance)

	mux := &http.ServeMux{}
	mux.Handle("/save", saveImageHandlerInstance)
	mux.Handle("/get", getImageHandlerInstance)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerConfig.Port),
		Handler: mux,
	}

	go func() {
		signalCh := make(chan os.Signal)
		signal.Notify(signalCh, os.Interrupt)

		select {
		case <-signalCh:
			logger.Info("stopping server gracefully")
			_ = server.Shutdown(ctx)
		case <-ctx.Done():
		}
	}()

	logger.Info(fmt.Sprintf("server starting on port %d", cfg.ServerConfig.Port))
	if err = server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server stopped with error", "error", err)
		}
	}

	logger.Info("server stopped")

	return 0
}
