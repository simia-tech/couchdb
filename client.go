package couchdb

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/simia-tech/errx"
)

// Client implements a simple couch db client.
type Client struct {
	BaseURL  string
	Username string
	Password string

	httpClient http.Client
}

// Database returns a reference to the database with the provided name.
func (c *Client) Database(name string) *DatabaseRef {
	return &DatabaseRef{
		client: c,
		name:   name,
	}
}

func (c *Client) do(request *http.Request) (*http.Response, error) {
	return c.httpClient.Do(request)
}

func (c *Client) requestFor(ctx context.Context, body io.Reader, method string, parts ...string) (*http.Request, error) {
	request, err := http.NewRequest(method, c.urlFor(parts...), body)
	if err != nil {
		return nil, errx.Annotatef(err, "new request")
	}
	if ctx != nil {
		request = request.WithContext(ctx)
	}

	if c.Username != "" && c.Password != "" {
		request.SetBasicAuth(c.Username, c.Password)
	}

	return request, nil
}

func (c *Client) urlFor(parts ...string) string {
	baseURL := c.BaseURL
	if strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL[:len(baseURL)-1]
	}
	return strings.Join(append([]string{c.BaseURL}, parts...), "/")
}
