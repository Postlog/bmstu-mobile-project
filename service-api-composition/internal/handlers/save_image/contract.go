package save_image

import (
	"context"
)

type imageStorageClientInterface interface {
	Save(ctx context.Context, imageBytes []byte) (string, error)
}
