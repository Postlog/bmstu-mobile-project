package save

import (
	"context"

	"github.com/google/uuid"
)

type imageRepository interface {
	Save(ctx context.Context, imageBytes []byte) (uuid.UUID, error)
}
