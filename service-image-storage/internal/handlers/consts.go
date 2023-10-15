package handlers

import "net/http"

const (
	ContentTypeJSON = "application/json"
)

const (
	ErrorCodeBadRequest          = http.StatusBadRequest
	ErrorCodeImageTooLarge       = 601
	ErrorCodeImageNotInPNGFormat = 602

	ErrorCodeImageNotFound = http.StatusNotFound
	ErrorCodeInternal      = http.StatusInternalServerError
)
