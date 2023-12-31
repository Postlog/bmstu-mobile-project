package save_image

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime"

	"github.com/postlog/mobile-project/service-image-storage/internal/handlers"
)

type Handler struct {
	logger    *slog.Logger
	imageRepo imageRepository
}

const (
	HTTPMethod = http.MethodGet
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

	imagesCount, err := h.imageRepo.Count(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		h.logger.WarnContext(r.Context(), "handler /info: unexpected repository error", "error", err)

		return
	}

	w.Header().Set("Content-Type", handlers.ContentTypeJSON)

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(Response{
		ImagesCount: imagesCount,
		OS:          runtime.GOOS,
	})
	if err != nil {
		h.logger.ErrorContext(r.Context(), "handler /info: error encoding response", "error", err)
	}
}
