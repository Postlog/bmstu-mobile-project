package handlers

type GetScaleResultRequest struct {
	TaskID string `json:"taskId"`
}

type GetScaleResultResponse struct {
	Result *GetScaleResultResponseResult `json:"result"`
	Error  *ResponseError                `json:"error"`
}

type GetScaleResultResponseResult struct {
	TaskID        string  `json:"taskId"`
	OriginImageID string  `json:"originImageId"`
	ScaleFactor   int     `json:"scaleFactor"`
	ImageID       *string `json:"imageId"`
	ScaleError    *string `json:"scaleError"`
}

type CreateScaleTaskRequest struct {
	ImageID     string `json:"imageId"`
	ScaleFactor int    `json:"scaleFactor"`
}

type CreateScaleTaskResponse struct {
	TaskID *string        `json:"taskId"`
	Error  *ResponseError `json:"error"`
}

type SaveResponse struct {
	ImageID *string        `json:"imageId"`
	Error   *ResponseError `json:"error"`
}

type ResponseError struct {
	Message string `json:"message"`
}
