package image_scaler

import (
	"context"
	"errors"
	"fmt"
	imageScalerClient "github.com/postlog/mobile-project/service-api-composition/internal/clients/image_scaler"
	"log/slog"

	"golang.org/x/sync/errgroup"

	scaleResultRepository "github.com/postlog/mobile-project/service-api-composition/internal/repository/scale_result"
	scaleTaskRepository "github.com/postlog/mobile-project/service-api-composition/internal/repository/scale_task"
)

type Service struct {
	logger *slog.Logger

	scaleResultRepository scaleResultRepositoryInterface
	scaleTaskRepository   scaleTaskRepositoryInterface
	imageScalerClient     imageScalerClientInterface
}

const (
	scaleTaskRepositoryBatchSize = 20
)

func New(
	logger *slog.Logger,
	scaleResultRepository scaleResultRepositoryInterface,
	scaleTaskRepository scaleTaskRepositoryInterface,
	imageScalerClient imageScalerClientInterface,
) *Service {
	return &Service{
		logger:                logger,
		scaleResultRepository: scaleResultRepository,
		scaleTaskRepository:   scaleTaskRepository,
		imageScalerClient:     imageScalerClient,
	}
}

func (s Service) Run(ctx context.Context) error {
	return s.scaleTaskRepository.Get(ctx, scaleTaskRepositoryBatchSize, s.processScaleTasks)
}

func (s Service) processScaleTasks(ctx context.Context, tasks []scaleTaskRepository.ScaleTask) error {
	eg, egCtx := errgroup.WithContext(ctx)

	s.logger.InfoContext(ctx, "new tasks consumed", "tasks", tasks)

	scaleResults := make([]scaleResultRepository.ScaleResult, len(tasks))

	for i, task := range tasks {
		t := task
		idx := i
		eg.Go(func() error {
			scaleResults[idx] = scaleResultRepository.ScaleResult{
				TaskID:          t.ID,
				OriginalImageID: t.ImageID,
				ScaleFactor:     t.ScaleFactor,
			}
			scaledImageID, err := s.imageScalerClient.ScaleImage(egCtx, t.ImageID, t.ScaleFactor)
			if err != nil {
				if errors.Is(err, imageScalerClient.ErrBadRequestValues) {
					tmp := "Некорректные значения для увеличения изображения"
					scaleResults[idx].ErrorText = &tmp

					return nil
				}

				s.logger.ErrorContext(ctx, "error scaling image", "task", t, "error", err)

				return fmt.Errorf("image-scaler ScaleImage: %w", err)
			}

			scaleResults[idx].ImageID = &scaledImageID

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	err := s.scaleResultRepository.Save(ctx, scaleResults)
	if err != nil {
		s.logger.ErrorContext(ctx, "error saving scale results", "results", scaleResults, "error", err)

		return fmt.Errorf("scale result repository Save: %w", err)
	}

	s.logger.InfoContext(ctx, "scale results saved", "results", scaleResults)

	return nil
}
