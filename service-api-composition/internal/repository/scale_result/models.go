package scale_result

type ScaleResult struct {
	TaskID        string
	OriginImageID string
	ScaleFactor   int
	ImageID       *string
	ErrorText     *string
}
