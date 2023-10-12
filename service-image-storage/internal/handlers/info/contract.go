package save_image

import (
	"context"
)

type imageRepository interface {
	Count(ctx context.Context) (int, error)
}
