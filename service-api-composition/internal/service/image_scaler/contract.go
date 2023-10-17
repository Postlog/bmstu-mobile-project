package image_scaler

import (
	"context"

	scaleResultRepository "github.com/postlog/mobile-project/service-api-composition/internal/repository/scale_result"
	scaleTaskRepository "github.com/postlog/mobile-project/service-api-composition/internal/repository/scale_task"
)

type scaleResultRepositoryInterface interface {
	Save(ctx context.Context, results []scaleResultRepository.ScaleResult) error
}

type scaleTaskRepositoryInterface interface {
	Get(ctx context.Context, batchSize int, processCallback func(context.Context, []scaleTaskRepository.ScaleTask) error) error
}

type imageScalerClientInterface interface {
	ScaleImage(ctx context.Context, imageID string, scaleFactor int) (string, error)
}
