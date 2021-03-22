package couchdb

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// HTTPClientStd implements a `HTTPClient` using go's standart library.
type HTTPClientStd struct {
	client *http.Client
}

// NewHTTPClientStd returns a new http client.
func NewHTTPClientStd(client *http.Client) *HTTPClientStd {
	if client == nil {
		client = http.DefaultClient
	}
	return &HTTPClientStd{client: client}
}

// Request performs a http request using the provided parameters.
func (c *HTTPClientStd) Request(ctx context.Context, method, url string, header http.Header) (io.ReadCloser, error) {
	request, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	request.Header = header

	response, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("do: %w", err)
	}

	return response.Body, nil
}
