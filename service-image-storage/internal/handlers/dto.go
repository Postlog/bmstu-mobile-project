package handlers

type GetRequest struct {
	ImageID string `json:"imageId"`
}

type GetResponse struct {
	EncodedImage *string        `json:"encodedImage,omitempty"`
	Error        *ResponseError `json:"error,omitempty"`
}

type SaveRequest struct {
	EncodedImage string `json:"encodedImage"`
}

type SaveResponse struct {
	ImageID *string        `json:"imageId,omitempty"`
	Error   *ResponseError `json:"error,omitempty"`
}

type ResponseError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}
