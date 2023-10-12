package save_image

import (
	"context"
	imageRepo "github.com/postlog/mobile-project/service-image-storage/internal/repository/image"

	"github.com/google/uuid"
)

type imageRepository interface {
	Get(ctx context.Context, id uuid.UUID) (imageRepo.Image, error)
}
