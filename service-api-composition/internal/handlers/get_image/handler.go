package save_image

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"

	imageStorageClient "github.com/postlog/mobile-project/service-api-composition/internal/clients/image_storage"
	"github.com/postlog/mobile-project/service-api-composition/internal/handlers"
)

type Handler struct {
	logger             *slog.Logger
	imageStorageClient imageStorageClientInterface
}

const (
	muxVarImageID = "imageId"
)

func New(logger *slog.Logger, imageStorageClient imageStorageClientInterface) *Handler {
	return &Handler{
		logger:             logger,
		imageStorageClient: imageStorageClient,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()

	vars := mux.Vars(r)
	imageID := vars[muxVarImageID]

	ctx := r.Context()

	imageBytes, err := h.imageStorageClient.Get(ctx, imageID)
	if err != nil {
		if errors.Is(err, imageStorageClient.ErrImageNotFound) {
			w.WriteHeader(http.StatusNotFound)

			h.logger.WarnContext(ctx, "handler get_image: image not found", "imageId", imageID)

			return
		}

		w.WriteHeader(http.StatusInternalServerError)

		h.logger.ErrorContext(ctx, "handler get_image: unexpected error from image-storage", "error", err, "imageId", imageID)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", handlers.ContentTypeImagePNG)

	_, _ = w.Write(imageBytes)
}
