package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/postlog/mobile-project/service-image-storage/internal/config"
	getImageHandler "github.com/postlog/mobile-project/service-image-storage/internal/handlers/get"
	infoHandler "github.com/postlog/mobile-project/service-image-storage/internal/handlers/info"
	saveImageHandler "github.com/postlog/mobile-project/service-image-storage/internal/handlers/save"
	imageRepository "github.com/postlog/mobile-project/service-image-storage/internal/repository/image"
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

	if err = os.MkdirAll(cfg.StorageConfig.FolderPath, os.ModePerm); err != nil {
		logger.ErrorContext(ctx, "error creating folder", "error", err, "folderPath", cfg.StorageConfig.FolderPath)
		return 1
	}

	logger.InfoContext(ctx, "initialized images directory", "directory", cfg.StorageConfig.FolderPath)

	imageRepo, err := imageRepository.New(cfg.StorageConfig.FolderPath)
	if err != nil {
		logger.ErrorContext(ctx, "error initializing image repository", "error", err)
		return 1
	}

	saveImageHandlerInstance := saveImageHandler.New(logger, imageRepo)
	infoHandlerInstance := infoHandler.New(logger, imageRepo)
	getImageHandlerInstance := getImageHandler.New(logger, imageRepo)

	mux := &http.ServeMux{}
	mux.Handle("/save", saveImageHandlerInstance)
	mux.Handle("/get", getImageHandlerInstance)
	mux.Handle("/info", infoHandlerInstance)

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
