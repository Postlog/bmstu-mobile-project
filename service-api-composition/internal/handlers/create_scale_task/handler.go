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
	HTTPMethod = http.MethodPost
)

const (
	errorMessageScaleFactorDoesntFitInRange = "Коэффициент масштабирования некорректен"
	errorMessageBadRequest                  = "Некорректный запрос"
	errorMessageInternalError               = "Непредвиденная ошибка, попробуйте позже"
)

func New(logger *slog.Logger, scaleTaskRepo scaleTaskRepositoryInterface) *Handler {
	return &Handler{
		logger:        logger,
		scaleTaskRepo: scaleTaskRepo,
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

	var req handlers.CreateScaleTaskRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		encodeErr := json.NewEncoder(w).Encode(handlers.SaveResponse{
			Error: &handlers.ResponseError{
				Message: errorMessageBadRequest,
			},
		})

		h.logger.WarnContext(r.Context(), "handler /createScaleTask: error decoding body", "error", err)

		if encodeErr != nil {
			h.logger.ErrorContext(r.Context(), "handler /createScaleTask: error encoding response", "error", encodeErr)
		}

		return
	}

	if req.ScaleFactor <= 0 {
		encodeErr := json.NewEncoder(w).Encode(handlers.SaveResponse{
			Error: &handlers.ResponseError{
				Message: errorMessageScaleFactorDoesntFitInRange,
			},
		})

		h.logger.WarnContext(r.Context(), "handler /createScaleTask: invalid scaleFactor", "scaleFactor", req.ScaleFactor)

		if encodeErr != nil {
			h.logger.ErrorContext(r.Context(), "handler /createScaleTask: error encoding response", "error", encodeErr)
		}

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
		encodeErr := json.NewEncoder(w).Encode(handlers.SaveResponse{
			Error: &handlers.ResponseError{
				Message: errorMessageInternalError,
			},
		})

		h.logger.ErrorContext(r.Context(), "handler /createScaleTask: unexpected error from scale task repository", "error", err)

		if encodeErr != nil {
			h.logger.ErrorContext(r.Context(), "handler /createScaleTask: error encoding response", "error", encodeErr)
		}

		return
	}

	err = json.NewEncoder(w).Encode(handlers.CreateScaleTaskResponse{
		TaskID: &taskID,
	})
	if err != nil {
		h.logger.ErrorContext(r.Context(), "handler /createScaleTask: error encoding response", "error", err)
	}
}
