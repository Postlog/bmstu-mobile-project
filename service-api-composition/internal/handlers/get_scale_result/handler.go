package create_scale_task

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/postlog/mobile-project/service-api-composition/internal/handlers"
	scaleResultRepository "github.com/postlog/mobile-project/service-api-composition/internal/repository/scale_result"
)

type Handler struct {
	logger                *slog.Logger
	scaleResultRepository scaleResultRepositoryInterface
}

const (
	HTTPMethod = http.MethodPost
)

const (
	errorMessageBadRequest    = "Некорректный запрос"
	errorMessageInternalError = "Непредвиденная ошибка, попробуйте позже"
)

func New(logger *slog.Logger, scaleResultRepository scaleResultRepositoryInterface) *Handler {
	return &Handler{
		logger:                logger,
		scaleResultRepository: scaleResultRepository,
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

	var req handlers.GetScaleResultRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		encodeErr := json.NewEncoder(w).Encode(handlers.GetScaleResultResponse{
			Error: &handlers.ResponseError{
				Message: errorMessageBadRequest,
			},
		})

		h.logger.WarnContext(r.Context(), "handler /getScaleResult: error decoding body", "error", err)

		if encodeErr != nil {
			h.logger.ErrorContext(r.Context(), "handler /getScaleResult: error encoding response", "error", encodeErr)
		}

		return
	}

	scaleResult, err := h.scaleResultRepository.Get(r.Context(), req.TaskID)
	if err != nil {
		if errors.Is(err, scaleResultRepository.ErrResultNotFound) {
			encodeErr := json.NewEncoder(w).Encode(handlers.GetScaleResultResponse{})

			if encodeErr != nil {
				h.logger.ErrorContext(r.Context(), "handler /getScaleResult: error encoding response", "error", encodeErr)
			}

			return
		}

		encodeErr := json.NewEncoder(w).Encode(handlers.GetScaleResultResponse{
			Error: &handlers.ResponseError{
				Message: errorMessageInternalError,
			},
		})

		h.logger.ErrorContext(r.Context(), "handler /getScaleResult: unexpected error from scale result repository", "error", err)

		if encodeErr != nil {
			h.logger.ErrorContext(r.Context(), "handler /getScaleResult: error encoding response", "error", encodeErr)
		}

		return
	}

	err = json.NewEncoder(w).Encode(handlers.GetScaleResultResponse{
		Result: &handlers.GetScaleResultResponseResult{
			TaskID:        req.TaskID,
			OriginImageID: scaleResult.OriginImageID,
			ScaleFactor:   scaleResult.ScaleFactor,
			ImageID:       scaleResult.ImageID,
			ScaleError:    scaleResult.ErrorText,
		},
	})
	if err != nil {
		h.logger.ErrorContext(r.Context(), "handler /getScaleResult: error encoding response", "error", err)
	}
}
