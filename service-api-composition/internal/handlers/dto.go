package handlers

type GetScaleResultResponse struct {
	Result *GetScaleResultResponseResult `json:"result,omitempty"`
	Error  *ResponseError                `json:"error,omitempty"`
}

type GetScaleResultResponseResult struct {
	TaskID          string  `json:"taskId"`
	OriginalImageID string  `json:"originalImageId"`
	ScaleFactor     int     `json:"scaleFactor"`
	ImageID         *string `json:"imageId,omitempty"`
	ScalingTime     *int    `json:"scalingResult,omitempty"`
	ScaleError      *string `json:"scaleError,omitempty"`
}

type CreateScaleTaskRequest struct {
	ImageID     string `json:"imageId"`
	ScaleFactor int    `json:"scaleFactor"`
}

type CreateScaleTaskResponse struct {
	TaskID *string        `json:"taskId,omitempty"`
	Error  *ResponseError `json:"error,omitempty"`
}

type SaveImageResponse struct {
	ImageID *string        `json:"imageId,omitempty"`
	Error   *ResponseError `json:"error,omitempty"`
}

type InfoResponse struct {
	OS string `json:"os"`
}

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
