package couchdb

import (
	"context"
	"io"
	"net/http"
)

// HTTPClient defines an interface of a http client.
type HTTPClient interface {
	Request(context.Context, string, string, http.Header) (io.ReadCloser, error)
}
