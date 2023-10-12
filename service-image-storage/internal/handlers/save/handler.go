package save

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

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
	wrongContentTypeMessage    = "unexpected content type"
	wrongBodyMessage           = "wrong body"
	imageTooLargeMessage       = "image too large"
	imageNotInPNGFormatMessage = "image not in png format"

	wrongContentTypeCode    = 1000
	wrongBodyCode           = 1001
	imageTooLargeCode       = 1002
	imageNotInPNGFormatCode = 1003
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

	contentType := r.Header.Get("Content-Type")
	if contentType != handlers.ContentTypeImagePNG {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Response{
			Error: &ResponseError{
				Code:    wrongContentTypeCode,
				Message: wrongContentTypeMessage,
			},
		})

		h.logger.WarnContext(r.Context(), "save_image: request with wrong content type", "contentType", contentType)

		return
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Response{
			Error: &ResponseError{
				Code:    wrongBodyCode,
				Message: wrongBodyMessage,
			},
		})

		h.logger.ErrorContext(r.Context(), "save_image: error reading body", "error", err)

		return
	}

	imageID, err := h.imageRepo.Save(r.Context(), bytes)
	if err != nil {
		if errors.Is(err, imageRepo.ErrImageTooLarge) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(Response{
				Error: &ResponseError{
					Code:    imageTooLargeCode,
					Message: imageTooLargeMessage,
				},
			})

			h.logger.WarnContext(r.Context(), "save_image: passed image is too large")

			return
		}
		if errors.Is(err, imageRepo.ErrImageNotInPNGFormat) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(Response{
				Error: &ResponseError{
					Code:    imageNotInPNGFormatCode,
					Message: imageNotInPNGFormatMessage,
				},
			})

			h.logger.WarnContext(r.Context(), "save_image: passed image not in png format")

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{})

		h.logger.WarnContext(r.Context(), "save_image: unexpected repository error", "error", err)

		return
	}

	w.WriteHeader(http.StatusOK)
	imageIDStr := imageID.String()
	_ = json.NewEncoder(w).Encode(Response{
		ImageID: &imageIDStr,
	})
}
