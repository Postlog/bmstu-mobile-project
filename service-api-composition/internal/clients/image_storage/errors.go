package image_storage

import "errors"

var (
	ErrImageNotInPNGFormat = errors.New("image not in png format")
	ErrImageTooLarge       = errors.New("image too large")
	ErrImageNotFound       = errors.New("image not found")
)
