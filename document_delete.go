package couchdb

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/simia-tech/errx"
)

// DocumentDeleteRequest defines the document fetch request.
type DocumentDeleteRequest struct {
	document *DocumentRef
	id       string
	revision string

	ctx context.Context
}

// WithContext adds a context to the request.
func (ddr *DocumentDeleteRequest) WithContext(ctx context.Context) *DocumentDeleteRequest {
	ddr.ctx = ctx
	return ddr
}

// Do performs the request.
func (ddr *DocumentDeleteRequest) Do() (*DocumentDeleteResponse, error) {
	request, err := ddr.document.database.client.requestFor(ddr.ctx, nil, http.MethodDelete, ddr.document.database.name, ddr.id)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", MimeTypeJSON)
	request.Header.Set("If-Match", ddr.revision)

	response, err := ddr.document.database.client.do(request)
	if err != nil {
		return nil, errx.Annotatef(err, "request [%s] [%s]", request.Method, request.URL)
	}
	defer response.Body.Close()

	r := &DocumentDeleteResponse{}
	if err := json.NewDecoder(response.Body).Decode(r); err != nil {
		return nil, errx.Annotatef(err, "json decode")
	}
	if err := evaluateResponseStatus(response.StatusCode, "",
		http.StatusOK, http.StatusAccepted,
		http.StatusBadRequest, http.StatusUnauthorized, http.StatusNotFound, http.StatusConflict); err != nil {
		return r, err
	}
	return r, nil
}

// DocumentDeleteResponse defines the document create response.
type DocumentDeleteResponse struct {
	OK       bool   `json:"ok"`
	ID       string `json:"id"`
	Revision string `json:"rev"`
}
