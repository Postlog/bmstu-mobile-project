package save_image

import (
	"errors"
	"log/slog"
	"net/http"

	imageStorage "github.com/postlog/mobile-project/service-api-composition/internal/clients/image_storage"
	"github.com/postlog/mobile-project/service-api-composition/internal/handlers"
)

type Handler struct {
	logger             *slog.Logger
	imageStorageClient imageStorageClient
}

const (
	HTTPMethod = http.MethodGet
)

const (
	imageIDKey = "id"
)

func New(logger *slog.Logger, imageStorageClient imageStorageClient) *Handler {
	return &Handler{
		logger:             logger,
		imageStorageClient: imageStorageClient,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()

	if r.Method != HTTPMethod {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	imageID := r.URL.Query().Get(imageIDKey)

	imageBytes, err := h.imageStorageClient.Get(r.Context(), imageID)
	if err != nil {
		if errors.Is(err, imageStorage.ErrImageNotFound) {
			w.WriteHeader(http.StatusNotFound)

			h.logger.WarnContext(r.Context(), "handler /get: image not exist", "imageId", imageID)

			return
		}

		w.WriteHeader(http.StatusInternalServerError)

		h.logger.ErrorContext(r.Context(), "handler /get: unexpected error from image-storage", "error", err, "imageId", imageID)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", handlers.ContentTypeImagePNG)
	_, err = w.Write(imageBytes)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "handler /get: error writing response", "error", err)
	}
}
