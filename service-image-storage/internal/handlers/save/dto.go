package save

type Response struct {
	ImageID *string        `json:"imageId,omitempty"`
	Error   *ResponseError `json:"error,omitempty"`
}

type ResponseError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}
