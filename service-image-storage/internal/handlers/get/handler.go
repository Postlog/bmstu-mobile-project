package save_image

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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

	var req handlers.GetRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		encodeErr := json.NewEncoder(w).Encode(handlers.GetResponse{
			Error: &handlers.ResponseError{
				Code:    handlers.ErrorCodeBadRequest,
				Message: "body must be valid json",
			},
		})

		h.logger.WarnContext(ctx, "handler /get: invalid json in request", "error", err)

		if encodeErr != nil {
			h.logger.ErrorContext(ctx, "handler /get: error encoding response", "error", encodeErr)
		}

		return
	}

	imageIDRaw := req.ImageID
	imageID, err := uuid.Parse(imageIDRaw)
	if err != nil {
		encodeErr := json.NewEncoder(w).Encode(handlers.GetResponse{
			Error: &handlers.ResponseError{
				Code:    handlers.ErrorCodeBadRequest,
				Message: "image id must be valid UUID",
			},
		})

		h.logger.WarnContext(ctx, "handler /get: invalid imageId format in request", "imageId", imageIDRaw, "error", err)

		if encodeErr != nil {
			h.logger.ErrorContext(ctx, "handler /get: error encoding response", "error", encodeErr)
		}

		return
	}

	image, err := h.imageRepo.Get(ctx, imageID)
	if err != nil {
		if errors.Is(err, imageRepo.ErrImageNotExist) {
			encodeErr := json.NewEncoder(w).Encode(handlers.GetResponse{
				Error: &handlers.ResponseError{
					Code:    handlers.ErrorCodeImageNotFound,
					Message: fmt.Sprintf("image with id '%s' not found", imageID.String()),
				},
			})

			h.logger.WarnContext(ctx, "handler /get: image with specified id not found", "imageId", imageID.String())

			if encodeErr != nil {
				h.logger.ErrorContext(ctx, "handler /get: error encoding response", "error", encodeErr)
			}

			return
		}

		encodeErr := json.NewEncoder(w).Encode(handlers.GetResponse{
			Error: &handlers.ResponseError{
				Code:    handlers.ErrorCodeInternal,
				Message: "unexpected internal error",
			},
		})

		h.logger.ErrorContext(
			ctx,
			"handler /get: unexpected error from image repository",
			"imageId", imageID.String(),
			"error", err,
		)

		if encodeErr != nil {
			h.logger.ErrorContext(ctx, "handler /get: error encoding response", "error", encodeErr)
		}

		return
	}

	encodedImage := base64.StdEncoding.EncodeToString(image.Bytes)

	err = json.NewEncoder(w).Encode(handlers.GetResponse{
		EncodedImage: &encodedImage,
	})

	if err != nil {
		h.logger.ErrorContext(ctx, "handler /get: error encoding response", "error", err)
	}
}
