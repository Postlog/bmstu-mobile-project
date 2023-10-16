package handlers

import "net/http"

const (
	ContentTypeJSON     = "application/json"
	ContentTypeImagePNG = "image/png"
)

const (
	ErrorCodeBadRequest    = http.StatusBadRequest
	ErrorCodeNotFound      = http.StatusNotFound
	ErrorCodeInternalError = http.StatusInternalServerError
)

const (
	ErrorMessageBodyMustBeJSON = "Тело запроса должно быть в формате JSON"
	ErrorMessageBadRequest     = "Некорректный запрос"
	ErrorMessageInternal       = "Непредвиденная ошибка, попробуйте позже"
)
