package save

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	imageStorage "github.com/postlog/mobile-project/service-api-composition/internal/clients/image_storage"
	"github.com/postlog/mobile-project/service-api-composition/internal/handlers"
)

type Handler struct {
	logger    *slog.Logger
	imageRepo imageStorageClient
}

const (
	HTTPMethod = http.MethodGet
)

const (
	errorMessageWrongImageFormat = "Изображение должно быть в формате PNG"
	errorMessageUnreadableBody   = "Тело запроса некорректно"
	errorMessageImageTooLarge    = "Изображение слишком большое"
	errorMessageInternalError    = "Непредвиденная ошибка, попробуйте позже"
)

func New(logger *slog.Logger, imageRepo imageStorageClient) *Handler {
	return &Handler{
		logger:    logger,
		imageRepo: imageRepo,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()

	if r.Method != HTTPMethod {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	w.Header().Set("Content-Type", handlers.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)

	contentType := r.Header.Get("Content-Type")
	if contentType != handlers.ContentTypeImagePNG {
		encodeErr := json.NewEncoder(w).Encode(handlers.SaveResponse{
			Error: &handlers.ResponseError{
				Message: errorMessageWrongImageFormat,
			},
		})

		h.logger.WarnContext(r.Context(), "handler /save: request with wrong content type", "contentType", contentType)

		if encodeErr != nil {
			h.logger.ErrorContext(r.Context(), "handler /save: error encoding response", "error", encodeErr)
		}

		return
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		encodeErr := json.NewEncoder(w).Encode(handlers.SaveResponse{
			Error: &handlers.ResponseError{
				Message: errorMessageUnreadableBody,
			},
		})

		h.logger.WarnContext(r.Context(), "handler /save: error reading body", "error", err)

		if encodeErr != nil {
			h.logger.ErrorContext(r.Context(), "handler /save: error encoding response", "error", encodeErr)
		}

		return
	}

	imageID, err := h.imageRepo.Save(r.Context(), bytes)
	if err != nil {
		if errors.Is(err, imageStorage.ErrImageTooLarge) {
			encodeErr := json.NewEncoder(w).Encode(handlers.SaveResponse{
				Error: &handlers.ResponseError{
					Message: errorMessageImageTooLarge,
				},
			})

			h.logger.WarnContext(r.Context(), "handler /save: passed image is too large")

			if encodeErr != nil {
				h.logger.ErrorContext(r.Context(), "handler /save: error encoding response", "error", encodeErr)
			}

			return
		}
		if errors.Is(err, imageStorage.ErrImageNotInPNGFormat) {
			encodeErr := json.NewEncoder(w).Encode(handlers.SaveResponse{
				Error: &handlers.ResponseError{
					Message: errorMessageWrongImageFormat,
				},
			})

			h.logger.WarnContext(r.Context(), "handler /save: passed image not in png format")

			if encodeErr != nil {
				h.logger.ErrorContext(r.Context(), "handler /save: error encoding response", "error", encodeErr)
			}

			return
		}

		encodeErr := json.NewEncoder(w).Encode(handlers.SaveResponse{
			Error: &handlers.ResponseError{
				Message: errorMessageInternalError,
			},
		})

		h.logger.ErrorContext(r.Context(), "handler /save: unexpected error from image-storage", "error", err)

		if encodeErr != nil {
			h.logger.ErrorContext(r.Context(), "handler /save: error encoding response", "error", encodeErr)
		}

		return
	}

	err = json.NewEncoder(w).Encode(handlers.SaveResponse{
		ImageID: &imageID,
	})
	if err != nil {
		h.logger.ErrorContext(r.Context(), "handler /save: error encoding response", "error", err)
	}
}
