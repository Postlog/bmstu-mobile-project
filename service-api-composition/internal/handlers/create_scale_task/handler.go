package create_scale_task

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/postlog/mobile-project/service-api-composition/internal/handlers"
	scaleTaskRepository "github.com/postlog/mobile-project/service-api-composition/internal/repository/scale_task"
)

type Handler struct {
	logger        *slog.Logger
	scaleTaskRepo scaleTaskRepositoryInterface
}

const (
	errorMessageScaleFactorDoesntFitInRange = "Коэффициент масштабирования некорректен"
)

func New(logger *slog.Logger, scaleTaskRepo scaleTaskRepositoryInterface) *Handler {
	return &Handler{
		logger:        logger,
		scaleTaskRepo: scaleTaskRepo,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()

	w.Header().Set("Content-Type", handlers.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)

	var req handlers.CreateScaleTaskRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		_ = json.NewEncoder(w).Encode(handlers.SaveImageResponse{
			Error: &handlers.ResponseError{
				Code:    handlers.ErrorCodeBadRequest,
				Message: handlers.ErrorMessageBodyMustBeJSON,
			},
		})

		h.logger.WarnContext(r.Context(), "handler create_scale_task: error decoding body into JSON", "error", err)

		return
	}

	if req.ScaleFactor <= 0 {
		_ = json.NewEncoder(w).Encode(handlers.SaveImageResponse{
			Error: &handlers.ResponseError{
				Code:    handlers.ErrorCodeBadRequest,
				Message: errorMessageScaleFactorDoesntFitInRange,
			},
		})

		h.logger.WarnContext(r.Context(), "handler create_scale_task: invalid scaleFactor", "scaleFactor", req.ScaleFactor)

		return
	}

	taskID := uuid.New().String()
	err = h.scaleTaskRepo.Save(r.Context(), []scaleTaskRepository.ScaleTask{
		{
			ID:          taskID,
			ImageID:     req.ImageID,
			ScaleFactor: req.ScaleFactor,
		},
	})
	if err != nil {
		_ = json.NewEncoder(w).Encode(handlers.SaveImageResponse{
			Error: &handlers.ResponseError{
				Code:    handlers.ErrorCodeInternalError,
				Message: handlers.ErrorMessageInternal,
			},
		})

		h.logger.ErrorContext(r.Context(), "handler create_scale_task: unexpected error from scale task repository", "error", err)

		return
	}

	_ = json.NewEncoder(w).Encode(handlers.CreateScaleTaskResponse{
		TaskID: &taskID,
	})
}
