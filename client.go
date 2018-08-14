package couchdb

import (
	"context"
	"io"
	"net/http"
	"strings"

	"code.posteo.de/common/errx"
)

// Client implements a simple couch db client.
type Client struct {
	BaseURL  string
	Username string
	Password string

	httpClient http.Client
}

// CreateDatabase creates a database of the provided name.
func (c *Client) CreateDatabase(ctx context.Context, name string) (*Database, error) {
	request, err := c.requestFor(ctx, nil, http.MethodPut, name)
	if err != nil {
		return nil, err
	}

	response, err := c.do(request)
	if err != nil {
		return nil, errx.Annotatef(err, "request [%s] [%s]", request.Method, request.URL)
	}
	defer response.Body.Close()

	_, err = evaluateResponse(response)
	if errx.IsAlreadyExists(err) {
		return nil, errx.AlreadyExistsf("database [%s] already exists", name)
	}
	if err != nil {
		return nil, err
	}

	return &Database{
		client: c,
		name:   name,
	}, nil
}

// GetDatabase returns a database object with the provided parameters.
func (c *Client) GetDatabase(name string) *Database {
	return &Database{
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
	request = request.WithContext(ctx)

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
