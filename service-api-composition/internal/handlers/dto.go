package handlers

type SaveResponse struct {
	ImageID *string        `json:"imageId,omitempty"`
	Error   *ResponseError `json:"error,omitempty"`
}

type ResponseError struct {
	Message string `json:"message"`
}
