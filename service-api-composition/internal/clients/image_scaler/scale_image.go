package image_storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c Client) ScaleImage(ctx context.Context, imageID string, scaleFactor int) (ScalingResult, error) {
	encodedBody, err := json.Marshal(scaleImageRequest{
		ImageID:     imageID,
		ScaleFactor: scaleFactor,
	})
	if err != nil {
		return ScalingResult{}, fmt.Errorf("error encoding request: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", c.baseURL, "scale"), bytes.NewReader(encodedBody))
	if err != nil {
		return ScalingResult{}, fmt.Errorf("error building request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return ScalingResult{}, fmt.Errorf("error sending request: %w", err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		return ScalingResult{}, fmt.Errorf("unexpected response status code %d", response.StatusCode)
	}

	var resp scaleImageResponse
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return ScalingResult{}, fmt.Errorf("error decoding response body to json: %w", err)
	}

	if resp.Error != nil {
		switch resp.Error.Code {
		case errorCodeBadRequest:
			return ScalingResult{}, ErrBadRequestValues
		}

		return ScalingResult{}, fmt.Errorf("unexpected error %d: %s", resp.Error.Code, resp.Error.Message)
	}

	if resp.Result == nil {
		return ScalingResult{}, fmt.Errorf("scaling result is nil")
	}

	return ScalingResult{
		ScaledImageID: resp.Result.ScaledImageID,
		ScalingTime:   resp.Result.ScalingTime,
	}, nil
}
