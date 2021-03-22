package couchdb

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/simia-tech/couchdb/value"
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

	if err := c.requestJSON(ctx, http.MethodGet, "/", &r); err != nil {
		return r, err
	}

	return r, nil
}

// AllDatabases fetches a list of all databases.
func (c *Client) AllDatabases(ctx context.Context) ([]string, error) {
	r := []string{}

	if err := c.requestJSON(ctx, http.MethodGet, "/_all_dbs", &r); err != nil {
		return r, err
	}

	return r, nil
}

func (c *Client) requestJSON(ctx context.Context, method, path string, responseBody interface{}) error {
	header := http.Header{}
	header.Add("Accept", "application/json")

	rc, err := c.request(ctx, method, path, header)
	if err != nil {
		return err
	}
	defer rc.Close()

	r := io.Reader(rc)
	// b := bytes.Buffer{}
	// r = io.TeeReader(r, &b)

	if err := json.NewDecoder(r).Decode(responseBody); err != nil {
		return fmt.Errorf("json decode: %w", err)
	}

	// log.Printf("json repsonse:\n%s\n", b.String())

	return nil
}

func (c *Client) request(ctx context.Context, method, path string, header http.Header) (io.ReadCloser, error) {
	url := c.baseURL + path

	if c.username != "" && c.password != "" {
		header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.username+":"+c.password)))
	}

	r, err := c.httpClient.Request(ctx, method, url, header)
	if err != nil {
		return nil, err
	}

	return r, nil
}
