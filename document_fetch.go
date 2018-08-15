package couchdb

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/simia-tech/errx"
)

// DocumentFetchRequest defines the document fetch request.
type DocumentFetchRequest struct {
	document *DocumentRef
	id       string

	ctx      context.Context
	revision string
}

// WithContext adds a context to the request.
func (dfr *DocumentFetchRequest) WithContext(ctx context.Context) *DocumentFetchRequest {
	dfr.ctx = ctx
	return dfr
}

// WithRevision adds a revision to the request.
func (dfr *DocumentFetchRequest) WithRevision(revision string) *DocumentFetchRequest {
	dfr.revision = revision
	return dfr
}

// Do performs the request.
func (dfr *DocumentFetchRequest) Do() (*DocumentFetchResponse, error) {
	request, err := dfr.document.database.client.requestFor(dfr.ctx, nil, http.MethodGet, dfr.document.database.name, dfr.id)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", MimeTypeJSON)
	request.Header.Set("Accept", MimeTypeJSON)
	if dfr.revision != "" {
		request.Header.Set("If-None-Match", dfr.revision)
	}

	response, err := dfr.document.database.client.do(request)
	if err != nil {
		return nil, errx.Annotatef(err, "request [%s] [%s]", request.Method, request.URL)
	}
	defer response.Body.Close()

	r := &DocumentFetchResponse{}
	if err := json.NewDecoder(response.Body).Decode(&r.Document); err != nil {
		return nil, errx.Annotatef(err, "json decode")
	}
	if err := evaluateResponseStatus(response.StatusCode, "",
		http.StatusOK,
		http.StatusNotModified,
		http.StatusBadRequest, http.StatusUnauthorized, http.StatusNotFound); err != nil {
		return r, err
	}
	return r, nil
}

// DocumentFetchResponse defines the document create response.
type DocumentFetchResponse struct {
	Document Document
}
