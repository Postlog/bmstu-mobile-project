package save_image

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/postlog/mobile-project/service-image-storage/internal/handlers"
	imageRepo "github.com/postlog/mobile-project/service-image-storage/internal/repository/image"
)

type Handler struct {
	logger    *slog.Logger
	imageRepo imageRepository
}

const (
	HTTPMethod = http.MethodGet
)

const (
	imageIDKey = "id"
)

func New(logger *slog.Logger, imageRepo imageRepository) *Handler {
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

	imageIDRaw := r.URL.Query().Get(imageIDKey)
	imageID, err := uuid.Parse(imageIDRaw)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		h.logger.WarnContext(r.Context(), "get_image: imageId in incorrect format", "error", err, "imageId", imageIDRaw)

		return
	}

	image, err := h.imageRepo.Get(r.Context(), imageID)
	if err != nil {
		if errors.Is(err, imageRepo.ErrImageNotExist) {
			w.WriteHeader(http.StatusNotFound)

			h.logger.WarnContext(r.Context(), "get_image: image not exist", "imageId", imageIDRaw)

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})

		h.logger.WarnContext(r.Context(), "get_image: unexpected repository error", "error", err, "imageId", imageIDRaw)

		return
	}

	w.Header().Set("Content-Type", handlers.ContentTypeImagePNG)
	_, _ = w.Write(image.Bytes)
}
