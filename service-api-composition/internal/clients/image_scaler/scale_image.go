package image_storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (c Client) ScaleImage(ctx context.Context, imageID string, scaleFactor int) (string, error) {
	time.Sleep(time.Second*3 + time.Second*time.Duration(rand.Intn(4)))

	if rand.Intn(100) > 70 {
		return "", fmt.Errorf("test error")
	}

	return uuid.New().String(), nil
}

func (c Client) scaleImageReal(ctx context.Context, imageID string, scaleFactor int) (string, error) {
	encodedBody, err := json.Marshal(scaleImageRequest{
		ImageID:     imageID,
		ScaleFactor: scaleFactor,
	})
	if err != nil {
		return "", fmt.Errorf("error encoding request: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", c.baseURL, "scaleImage"), bytes.NewReader(encodedBody))
	if err != nil {
		return "", fmt.Errorf("error building request: %w", err)
	}

	response, err := c.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response status code %d", response.StatusCode)
	}

	var resp scaleImageResponse
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return "", fmt.Errorf("error decoding response body to json: %w", err)
	}

	if resp.Error != nil {
		switch resp.Error.Code {
		}

		return "", fmt.Errorf("unexpected error %d: %s", resp.Error.Code, resp.Error.Message)
	}

	if resp.ScaledImageID == nil {
		return "", fmt.Errorf("scaled image id is nil")
	}

	return *resp.ScaledImageID, nil
}
