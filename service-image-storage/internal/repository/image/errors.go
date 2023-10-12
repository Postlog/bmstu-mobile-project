package image

import "errors"

var (
	ErrImageNotInPNGFormat = errors.New("image not in png format")
	ErrImageTooLarge       = errors.New("image too large")
	ErrImageNotExist       = errors.New("image with specified id not exist")
)
