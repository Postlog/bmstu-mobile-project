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

const (
	wrongContentTypeMessage = "unexpected content type"
)

func New(logger *slog.Logger, imageRepo imageRepository) *Handler {
	return &Handler{
		logger:    logger,
		imageRepo: imageRepo,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()

	w.Header().Set("Content-Type", handlers.ContentTypeJSON)

	if r.Method != HTTPMethod {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	imagesCount, err := h.imageRepo.Count(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})

		h.logger.WarnContext(r.Context(), "info: unexpected repository error", "error", err)

		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(Response{
		ImagesCount: imagesCount,
		OS:          runtime.GOOS,
	})
}
