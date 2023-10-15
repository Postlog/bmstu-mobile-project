package save

import (
	"context"
)

type imageStorageClient interface {
	Save(ctx context.Context, imageBytes []byte) (string, error)
}
