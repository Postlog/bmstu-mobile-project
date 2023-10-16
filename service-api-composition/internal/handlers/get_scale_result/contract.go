package create_scale_task

import (
	"context"

	scaleResultRepository "github.com/postlog/mobile-project/service-api-composition/internal/repository/scale_result"
)

type scaleResultRepositoryInterface interface {
	Get(ctx context.Context, taskID string) (scaleResultRepository.ScaleResult, error)
}
