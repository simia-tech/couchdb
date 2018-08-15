package couchdb

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/simia-tech/errx"
)

// DocumentCreateRequest defines the document create request.
type DocumentCreateRequest struct {
	document *DocumentRef
	doc      Document

	ctx context.Context
}

// WithContext adds a context to the request.
func (dcr *DocumentCreateRequest) WithContext(ctx context.Context) *DocumentCreateRequest {
	dcr.ctx = ctx
	return dcr
}

// Do performs the request.
func (dcr *DocumentCreateRequest) Do() (*DocumentCreateResponse, error) {
	buffer := bytes.Buffer{}
	if err := json.NewEncoder(&buffer).Encode(dcr.doc); err != nil {
		return nil, errx.Annotatef(err, "json encode")
	}

	request, err := dcr.document.database.client.requestFor(dcr.ctx, &buffer, http.MethodPost, dcr.document.database.name)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", MimeTypeJSON)
	request.Header.Set("Accept", MimeTypeJSON)

	response, err := dcr.document.database.client.do(request)
	if err != nil {
		return nil, errx.Annotatef(err, "request [%s] [%s]", request.Method, request.URL)
	}
	defer response.Body.Close()

	r := &DocumentCreateResponse{}
	if err := json.NewDecoder(response.Body).Decode(r); err != nil {
		return nil, errx.Annotatef(err, "json decode")
	}
	if err := evaluateResponseStatus(response.StatusCode, "",
		http.StatusCreated, http.StatusAccepted,
		http.StatusBadRequest, http.StatusUnauthorized, http.StatusNotFound, http.StatusConflict); err != nil {
		return r, err
	}
	return r, nil
}

// DocumentCreateResponse defines the document create response.
type DocumentCreateResponse struct {
	OK       bool   `json:"ok"`
	ID       string `json:"id"`
	Revision string `json:"rev"`
}
