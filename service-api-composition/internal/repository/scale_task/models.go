package scale_task

type ScaleTask struct {
	ID          string `json:"id"`
	ImageID     string `json:"imageId"`
	ScaleFactor int    `json:"scaleFactor"`
}
