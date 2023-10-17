package scale_result

type ScaleResult struct {
	TaskID          string
	OriginalImageID string
	ScaleFactor     int
	ImageID         *string
	ErrorText       *string
}
