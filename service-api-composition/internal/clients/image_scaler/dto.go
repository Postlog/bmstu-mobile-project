package image_storage

type scaleImageRequest struct {
	ImageID     string `json:"imageId"`
	ScaleFactor int    `json:"scaleFactor"`
}

type scaleImageResponse struct {
	Result *scaleImageResponseScalingResult `json:"result"`
	Error  *responseError                   `json:"error"`
}

type scaleImageResponseScalingResult struct {
	ScaledImageID string `json:"scaledImageId"`
	ScalingTime   int    `json:"scalingTime"`
}

type responseError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

/*
{
	"result": {
		"scaledImageId": "",
		"scalingTime": ""
	}
}
*/
