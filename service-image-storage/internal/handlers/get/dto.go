package save_image

type Request struct {
	ImageID string `json:"imageId"`
}

type Response struct {
	Error *ResponseError `json:"error,omitempty"`
}

type ResponseError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}
