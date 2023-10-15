package image_storage

type getRequest struct {
	ImageID string `json:"imageId"`
}

type getResponse struct {
	EncodedImage *string        `json:"encodedImage,omitempty"`
	Error        *responseError `json:"error,omitempty"`
}

type saveRequest struct {
	EncodedImage string `json:"encodedImage"`
}

type saveResponse struct {
	ImageID *string        `json:"imageId,omitempty"`
	Error   *responseError `json:"error,omitempty"`
}

type responseError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}
