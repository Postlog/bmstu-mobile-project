package save_image

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime"

	"github.com/postlog/mobile-project/service-api-composition/internal/handlers"
)

type Handler struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()

	w.Header().Set("Content-Type", handlers.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(handlers.InfoResponse{
		OS: runtime.GOOS,
	})
}
