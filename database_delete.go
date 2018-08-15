package couchdb

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/simia-tech/errx"
)

// DatabaseDeleteRequest defines the database delete request.
type DatabaseDeleteRequest struct {
	database *DatabaseRef

	ctx context.Context
}

// WithContext adds a context to the request.
func (ddr *DatabaseDeleteRequest) WithContext(ctx context.Context) *DatabaseDeleteRequest {
	ddr.ctx = ctx
	return ddr
}

// Do performs the request.
func (ddr *DatabaseDeleteRequest) Do() (*DatabaseDeleteResponse, error) {
	request, err := ddr.database.client.requestFor(ddr.ctx, nil, http.MethodDelete, ddr.database.name)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", MimeTypeJSON)

	response, err := ddr.database.client.do(request)
	if err != nil {
		return nil, errx.Annotatef(err, "request %s %s", request.Method, request.URL)
	}
	defer response.Body.Close()

	r := &DatabaseDeleteResponse{}
	if err := json.NewDecoder(response.Body).Decode(r); err != nil {
		return nil, errx.Annotatef(err, "decode json")
	}
	if err := evaluateResponseStatus(response.StatusCode, "",
		http.StatusOK,
		http.StatusBadRequest, http.StatusUnauthorized, http.StatusNotFound); err != nil {
		return r, err
	}
	return r, nil
}

// DatabaseDeleteResponse defines the database delete response.
type DatabaseDeleteResponse struct {
	OK bool `json:"ok"`
}
