package save_image

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	imageStorageClient "github.com/postlog/mobile-project/service-api-composition/internal/clients/image_storage"
	"github.com/postlog/mobile-project/service-api-composition/internal/handlers"
)

type Handler struct {
	logger             *slog.Logger
	imageStorageClient imageStorageClientInterface
}

const (
	errorMessageWrongImageFormat = "Изображение должно быть в формате PNG"
	errorMessageImageTooLarge    = "Изображение слишком большое"
)

func New(logger *slog.Logger, imageRepo imageStorageClientInterface) *Handler {
	return &Handler{
		logger:             logger,
		imageStorageClient: imageRepo,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()

	w.Header().Set("Content-Type", handlers.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)

	ctx := r.Context()

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		_ = json.NewEncoder(w).Encode(handlers.SaveImageResponse{
			Error: &handlers.ResponseError{
				Code:    handlers.ErrorCodeBadRequest,
				Message: handlers.ErrorMessageBadRequest,
			},
		})

		h.logger.WarnContext(ctx, "handler save_image: error reading body", "error", err)

		return
	}

	imageID, err := h.imageStorageClient.Save(ctx, bytes)
	if err != nil {
		errorCode := handlers.ErrorCodeInternalError
		message := handlers.ErrorMessageInternal
		if errors.Is(err, imageStorageClient.ErrImageTooLarge) {
			errorCode = handlers.ErrorCodeBadRequest
			message = errorMessageImageTooLarge

			h.logger.WarnContext(ctx, "handler save_image: passed image is too large")
		} else if errors.Is(err, imageStorageClient.ErrImageNotInPNGFormat) {
			errorCode = handlers.ErrorCodeBadRequest
			message = errorMessageWrongImageFormat

			h.logger.WarnContext(ctx, "handler save_image: passed image not in png format")
		} else {
			h.logger.ErrorContext(ctx, "handler save_image: unexpected error from image-storage", "error", err)
		}

		_ = json.NewEncoder(w).Encode(handlers.SaveImageResponse{
			Error: &handlers.ResponseError{
				Code:    errorCode,
				Message: message,
			},
		})

		return
	}

	_ = json.NewEncoder(w).Encode(handlers.SaveImageResponse{
		ImageID: &imageID,
	})
}
