package image_storage

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c Client) Save(ctx context.Context, imageBytes []byte) (string, error) {
	encodedImage := base64.StdEncoding.EncodeToString(imageBytes)
	encodedBody, err := json.Marshal(saveRequest{
		EncodedImage: encodedImage,
	})
	if err != nil {
		return "", fmt.Errorf("error encoding request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", c.baseURL, "save"), bytes.NewReader(encodedBody))
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

	var resp saveResponse
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return "", fmt.Errorf("error decoding response body to json: %w", err)
	}

	if resp.Error != nil {
		switch resp.Error.Code {
		case errorCodeBadRequest:
			return "", fmt.Errorf("bad request: %s", resp.Error.Message)
		case errorCodeImageTooLarge:
			return "", ErrImageTooLarge
		case errorCodeImageNotInPNGFormat:
			return "", ErrImageNotInPNGFormat
		}

		return "", fmt.Errorf("unexpected error %d: %s", resp.Error.Code, resp.Error.Message)
	}

	if resp.ImageID == nil {
		return "", fmt.Errorf("imageID is nil")
	}

	return *resp.ImageID, nil
}
