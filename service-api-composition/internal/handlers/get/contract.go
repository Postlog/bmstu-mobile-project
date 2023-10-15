package save_image

import (
	"context"
)

type imageStorageClient interface {
	Get(ctx context.Context, id string) ([]byte, error)
}
