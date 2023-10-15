package save

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/postlog/mobile-project/service-image-storage/internal/handlers"
	imageRepo "github.com/postlog/mobile-project/service-image-storage/internal/repository/image"
)

type Handler struct {
	logger    *slog.Logger
	imageRepo imageRepository
}

const (
	HTTPMethod = http.MethodPost
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

	w.Header().Set("Content-Type", handlers.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)

	ctx := r.Context()

	var req handlers.SaveRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		encodeErr := json.NewEncoder(w).Encode(handlers.SaveResponse{
			Error: &handlers.ResponseError{
				Code:    handlers.ErrorCodeBadRequest,
				Message: "body must be valid json",
			},
		})

		h.logger.WarnContext(ctx, "handler /save: invalid json in request", "error", err)

		if encodeErr != nil {
			h.logger.ErrorContext(ctx, "handler /save: error encoding response", "error", encodeErr)
		}

		return
	}

	imageBytes, err := io.ReadAll(base64.NewDecoder(base64.StdEncoding, strings.NewReader(req.EncodedImage)))
	if err != nil {
		encodeErr := json.NewEncoder(w).Encode(handlers.SaveResponse{
			Error: &handlers.ResponseError{
				Code:    handlers.ErrorCodeBadRequest,
				Message: "encoded image must be valid base64",
			},
		})

		h.logger.WarnContext(ctx, "handler /save: base64 in request", "error", err)

		if encodeErr != nil {
			h.logger.ErrorContext(ctx, "handler /save: error encoding response", "error", encodeErr)
		}

		return
	}

	imageID, err := h.imageRepo.Save(ctx, imageBytes)
	if err != nil {
		if errors.Is(err, imageRepo.ErrImageTooLarge) {
			encodeErr := json.NewEncoder(w).Encode(handlers.SaveResponse{
				Error: &handlers.ResponseError{
					Code:    handlers.ErrorCodeImageTooLarge,
					Message: "provided image is too large",
				},
			})

			h.logger.WarnContext(ctx, "handler /save: provided image is too large")

			if encodeErr != nil {
				h.logger.ErrorContext(ctx, "handler /save: error encoding response", "error", encodeErr)
			}

			return
		}

		if errors.Is(err, imageRepo.ErrImageNotInPNGFormat) {
			encodeErr := json.NewEncoder(w).Encode(handlers.SaveResponse{
				Error: &handlers.ResponseError{
					Code:    handlers.ErrorCodeImageNotInPNGFormat,
					Message: "provided image is not in PNG format",
				},
			})

			h.logger.WarnContext(ctx, "handler /save: provided image is not in PNG format")

			if encodeErr != nil {
				h.logger.ErrorContext(ctx, "handler /save: error encoding response", "error", encodeErr)
			}

			return
		}

		encodeErr := json.NewEncoder(w).Encode(handlers.SaveResponse{
			Error: &handlers.ResponseError{
				Code:    handlers.ErrorCodeInternal,
				Message: "unexpected internal error",
			},
		})

		h.logger.ErrorContext(
			ctx,
			"handler /save: unexpected error from image repository",
			"error", err,
		)

		if encodeErr != nil {
			h.logger.ErrorContext(ctx, "handler /save: error encoding response", "error", encodeErr)
		}

		return
	}

	imageIDStr := imageID.String()
	err = json.NewEncoder(w).Encode(handlers.SaveResponse{
		ImageID: &imageIDStr,
	})

	if err != nil {
		h.logger.ErrorContext(ctx, "handler /save: error encoding response", "error", err)
	}
}
