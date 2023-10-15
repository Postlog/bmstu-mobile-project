package image_storage

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c Client) Get(ctx context.Context, id string) ([]byte, error) {
	encodedBody, err := json.Marshal(getRequest{
		ImageID: id,
	})
	if err != nil {
		return nil, fmt.Errorf("error encoding request: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", c.baseURL, "get"), bytes.NewReader(encodedBody))
	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}

	response, err := c.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status code %d", response.StatusCode)
	}

	var resp getResponse
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return nil, fmt.Errorf("error decoding response body to json: %w", err)
	}

	if resp.Error != nil {
		switch resp.Error.Code {
		case errorCodeBadRequest:
			return nil, fmt.Errorf("bad request: %s", resp.Error.Message)
		case errorCodeImageNotFound:
			return nil, ErrImageNotFound
		}

		return nil, fmt.Errorf("unexpected error %d: %s", resp.Error.Code, resp.Error.Message)
	}

	if resp.EncodedImage == nil {
		return nil, fmt.Errorf("encoded image is nil")
	}

	imageBytes, err := io.ReadAll(base64.NewDecoder(base64.StdEncoding, strings.NewReader(*resp.EncodedImage)))
	if err != nil {
		return nil, fmt.Errorf("error decoding image from base64: %w", err)
	}

	return imageBytes, nil
}
