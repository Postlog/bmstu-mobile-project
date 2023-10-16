package create_scale_task

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
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
	errorMessageResultNotFound = "Результат не найден"
	errorMessageInternalError  = "Непредвиденная ошибка, попробуйте позже"

	muxVarTaskID = "taskId"
)

func New(logger *slog.Logger, scaleResultRepository scaleResultRepositoryInterface) *Handler {
	return &Handler{
		logger:                logger,
		scaleResultRepository: scaleResultRepository,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()

	w.Header().Set("Content-Type", handlers.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	taskID := vars[muxVarTaskID]

	ctx := r.Context()

	scaleResult, err := h.scaleResultRepository.Get(ctx, taskID)
	if err != nil {
		if errors.Is(err, scaleResultRepository.ErrResultNotFound) {
			_ = json.NewEncoder(w).Encode(handlers.GetScaleResultResponse{
				Error: &handlers.ResponseError{
					Code:    handlers.ErrorCodeNotFound,
					Message: errorMessageResultNotFound,
				},
			})

			h.logger.WarnContext(ctx, "handler get_scale_result: scale result not found", "taskId", taskID, "error", err)

			return
		}

		_ = json.NewEncoder(w).Encode(handlers.GetScaleResultResponse{
			Error: &handlers.ResponseError{
				Code:    handlers.ErrorCodeInternalError,
				Message: handlers.ErrorMessageInternal,
			},
		})

		h.logger.ErrorContext(ctx, "handler get_scale_result: unexpected error from scale result repository", "error", err)

		return
	}

	_ = json.NewEncoder(w).Encode(handlers.GetScaleResultResponse{
		Result: &handlers.GetScaleResultResponseResult{
			TaskID:        taskID,
			OriginImageID: scaleResult.OriginImageID,
			ScaleFactor:   scaleResult.ScaleFactor,
			ImageID:       scaleResult.ImageID,
			ScaleError:    scaleResult.ErrorText,
		},
	})
}
