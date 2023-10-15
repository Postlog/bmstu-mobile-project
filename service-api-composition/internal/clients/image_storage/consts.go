package image_storage

import "net/http"

const (
	errorCodeBadRequest          = http.StatusBadRequest
	errorCodeImageTooLarge       = 601
	errorCodeImageNotInPNGFormat = 602

	errorCodeImageNotFound = http.StatusNotFound
)
