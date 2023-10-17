package image_storage

type scaleImageRequest struct {
	ImageID     string `json:"imageId"`
	ScaleFactor int    `json:"scaleFactor"`
}

type scaleImageResponse struct {
	ScaledImageID *string        `json:"scaledImageId"`
	Error         *responseError `json:"error"`
}

type responseError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}
