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

const (
	HTTPMethod = http.MethodGet
)

func New(logger *slog.Logger) *Handler {
	return &Handler{
		logger: logger,
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
	err := json.NewEncoder(w).Encode(Response{
		OS: runtime.GOOS,
	})
	if err != nil {
		h.logger.ErrorContext(r.Context(), "handler /info: error encoding response", "error", err)
	}
}
