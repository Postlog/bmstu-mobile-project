package save_image

import (
	"context"
)

type imageStorageClientInterface interface {
	Get(ctx context.Context, id string) ([]byte, error)
}
