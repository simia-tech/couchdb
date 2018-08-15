package couchdb

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/simia-tech/errx"
)

// DatabaseCreateRequest defines the database create request.
type DatabaseCreateRequest struct {
	database *DatabaseRef

	ctx context.Context
}

// WithContext adds a context to the request.
func (dcr *DatabaseCreateRequest) WithContext(ctx context.Context) *DatabaseCreateRequest {
	dcr.ctx = ctx
	return dcr
}

// Do performs the request.
func (dcr *DatabaseCreateRequest) Do() (*DatabaseCreateResponse, error) {
	request, err := dcr.database.client.requestFor(dcr.ctx, nil, http.MethodPut, dcr.database.name)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", MimeTypeJSON)

	response, err := dcr.database.client.do(request)
	if err != nil {
		return nil, errx.Annotatef(err, "request %s %s", request.Method, request.URL)
	}
	defer response.Body.Close()

	r := &DatabaseCreateResponse{}
	if err := json.NewDecoder(response.Body).Decode(r); err != nil {
		return nil, errx.Annotatef(err, "decode json")
	}
	if err := evaluateResponseStatus(response.StatusCode, r.Reason,
		http.StatusCreated,
		http.StatusBadRequest, http.StatusUnauthorized, http.StatusPreconditionFailed); err != nil {
		return r, err
	}
	return r, nil
}

// DatabaseCreateResponse efines the database create response.
type DatabaseCreateResponse struct {
	OK     bool   `json:"ok"`
	Error  string `json:"error"`
	Reason string `json:"reason"`
}
