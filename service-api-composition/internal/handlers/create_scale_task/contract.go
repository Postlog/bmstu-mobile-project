package create_scale_task

import (
	"context"

	scaleTaskRepository "github.com/postlog/mobile-project/service-api-composition/internal/repository/scale_task"
)

type scaleTaskRepositoryInterface interface {
	Save(ctx context.Context, tasks []scaleTaskRepository.ScaleTask) error
}
