package couchdb

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/simia-tech/couchdb/value"
)

// Various errors.
var (
	ErrBadRequest = errors.New("bad request")
	ErrConflict   = errors.New("conflict")
)

// Client implements a simple couch db client.
type Client struct {
	baseURL  string
	username string
	password string

	httpClient HTTPClient
}

// NewClient returns a new client configured with the provided options.
func NewClient(baseURL string, options ...ClientOption) (*Client, error) {
	c := &Client{
		baseURL:    baseURL,
		httpClient: NewHTTPClientStd(nil),
	}
	for _, o := range options {
		if err := o(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// InstanceInfo fetches some basic infos from the couchdb instance.
func (c *Client) InstanceInfo(ctx context.Context) (value.InstanceInfo, error) {
	r := value.InstanceInfo{}

	if err := c.requestJSON(ctx, http.MethodGet, "/", nil, nil, &r); err != nil {
		return r, err
	}

	return r, nil
}

// AllDatabases fetches a list of all databases.
func (c *Client) AllDatabases(ctx context.Context) ([]string, error) {
	r := []string{}

	if err := c.requestJSON(ctx, http.MethodGet, "/_all_dbs", nil, nil, &r); err != nil {
		return r, err
	}

	return r, nil
}

func (c *Client) requestJSON(
	ctx context.Context,
	method,
	path string,
	header http.Header,
	body,
	responseBody interface{},
) error {
	if header == nil {
		header = http.Header{}
	}
	header.Add("Accept", "application/json")

	bodyReader := io.Reader(nil)
	if body != nil {
		buffer := &bytes.Buffer{}
		if err := json.NewEncoder(buffer).Encode(body); err != nil {
			return fmt.Errorf("json encode: %w", err)
		}
		bodyReader = buffer
		header.Add("Content-Type", "application/json")
	}

	statusCode, _, responseReader, err := c.request(ctx, method, path, header, bodyReader)
	if err != nil {
		return err
	}
	defer responseReader.Close()

	reader := io.Reader(responseReader)
	// b := bytes.Buffer{}
	// reader = io.TeeReader(reader, &b)

	if err := checkJSONError(statusCode, reader); err != nil {
		return err
	}

	if err := json.NewDecoder(reader).Decode(responseBody); err != nil {
		return fmt.Errorf("json decode: %w", err)
	}

	// log.Printf("json repsonse:\n%s\n", b.String())

	return nil
}

func (c *Client) request(
	ctx context.Context,
	method,
	path string,
	header http.Header,
	body io.Reader,
) (int, http.Header, io.ReadCloser, error) {
	url := c.baseURL + path

	if c.username != "" && c.password != "" {
		header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.username+":"+c.password)))
	}

	statusCode, responseHeader, responseReader, err := c.httpClient.Request(ctx, method, url, header, body)
	if err != nil {
		return 0, nil, nil, err
	}

	return statusCode, responseHeader, responseReader, nil
}

func checkJSONError(statusCode int, reader io.Reader) error {
	switch statusCode {
	case http.StatusBadRequest:
		return ErrBadRequest
	case http.StatusConflict:
		return ErrConflict
	}
	return nil
}
