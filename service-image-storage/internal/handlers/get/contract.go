package save_image

import (
	"context"

	"github.com/google/uuid"

	imageRepo "github.com/postlog/mobile-project/service-image-storage/internal/repository/image"
)

type imageRepository interface {
	Get(ctx context.Context, id uuid.UUID) (imageRepo.Image, error)
}
