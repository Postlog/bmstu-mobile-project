package image_storage

import (
	"net/http"
)

type Client struct {
	baseURL    string
	httpClient http.Client
}

func New(url string, c http.Client) *Client {
	return &Client{
		baseURL:    url,
		httpClient: c,
	}
}
