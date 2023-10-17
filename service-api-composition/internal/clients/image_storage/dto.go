package image_storage

type getRequest struct {
	ImageID string `json:"imageId"`
}

type getResponse struct {
	EncodedImage *string        `json:"encodedImage"`
	Error        *responseError `json:"error"`
}

type saveRequest struct {
	EncodedImage string `json:"encodedImage"`
}

type saveResponse struct {
	ImageID *string        `json:"imageId"`
	Error   *responseError `json:"error"`
}

type responseError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}
