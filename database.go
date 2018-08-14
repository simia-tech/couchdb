package couchdb

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"code.posteo.de/common/errx"
)

// Database implements a couchdb database.
type Database struct {
	client *Client
	name   string
}

// CreateDocument creates a new document in the couch db database and returns the id and revision.
func (d *Database) CreateDocument(ctx context.Context, document Document) (string, string, error) {
	buffer := bytes.Buffer{}
	if err := json.NewEncoder(&buffer).Encode(document); err != nil {
		return "", "", errx.Annotatef(err, "json encode")
	}

	request, err := d.client.requestFor(ctx, &buffer, http.MethodPost, d.name)
	if err != nil {
		return "", "", err
	}

	request.Header.Set("Content-Type", MimeTypeJSON)
	request.Header.Set("Accept", MimeTypeJSON)

	response, err := d.client.do(request)
	if err != nil {
		return "", "", errx.Annotatef(err, "request [%s] [%s]", request.Method, request.URL)
	}
	defer response.Body.Close()

	result, err := evaluateResponse(response)
	if err != nil {
		return "", "", err
	}

	return getStringFieldOr(result, "id", ""), getStringFieldOr(result, "rev", ""), nil
}

// FetchDocument fetches the document with the provided id and rev from the couchdb database and returns it. If
// the provided rev is an empty string, the latest revision is returned.
func (d *Database) FetchDocument(ctx context.Context, id, rev string) (Document, error) {
	request, err := d.client.requestFor(ctx, nil, http.MethodGet, d.name, id)
	if err != nil {
		return nil, err
	}

	if rev != "" {
		request.Header.Set("If-None-Match", rev)
	}
	request.Header.Set("Accept", MimeTypeJSON)

	response, err := d.client.do(request)
	if err != nil {
		return nil, errx.Annotatef(err, "request [%s] [%s]", request.Method, request.URL)
	}
	defer response.Body.Close()

	result, err := evaluateResponse(response)
	if errx.IsNotFound(err) {
		return nil, errx.NotFoundf("could not find document with id [%s] and revision [%s]", id, rev)
	}
	if err != nil {
		return nil, err
	}

	if m, ok := result.(map[string]interface{}); ok {
		return m, nil
	}

	log.Printf("unexpected body: %+v", result)
	return nil, errx.BadRequestf("unexpected response %+v", result)
}

// FetchDocumentMeta returns the latest revision and content length of the document at the provided id.
func (d *Database) FetchDocumentMeta(ctx context.Context, id string) (string, uint64, error) {
	request, err := d.client.requestFor(ctx, nil, http.MethodHead, d.name, id)
	if err != nil {
		return "", 0, err
	}

	request.Header.Set("Accept", MimeTypeJSON)

	response, err := d.client.do(request)
	if err != nil {
		return "", 0, errx.Annotatef(err, "request [%s] [%s]", request.Method, request.URL)
	}
	defer response.Body.Close()

	err = evaluateResponseStatus(response.StatusCode, nil)
	if errx.IsNotFound(err) {
		return "", 0, errx.NotFoundf("could not find document revision with id [%s]", id)
	}
	if err != nil {
		return "", 0, err
	}

	revision := strings.Trim(response.Header.Get("ETag"), "\"")
	contentLength, err := strconv.ParseUint(response.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return "", 0, errx.Annotatef(err, "parse int64 [%s]", response.Header.Get("Content-Length"))
	}
	return revision, contentLength, nil
}

// UpdateDocument updates the document at the provided id.
func (d *Database) UpdateDocument(ctx context.Context, id, rev string, document Document) (string, error) {
	buffer := bytes.Buffer{}
	if err := json.NewEncoder(&buffer).Encode(document); err != nil {
		return "", errx.Annotatef(err, "json encode")
	}

	request, err := d.client.requestFor(ctx, &buffer, http.MethodPut, d.name, id)
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", MimeTypeJSON)
	request.Header.Set("If-Match", rev)
	request.Header.Set("Accept", MimeTypeJSON)

	response, err := d.client.do(request)
	if err != nil {
		return "", errx.Annotatef(err, "request [%s] [%s]", request.Method, request.URL)
	}
	defer response.Body.Close()

	result, err := evaluateResponse(response)
	if err != nil {
		return "", err
	}

	return getStringFieldOr(result, "rev", ""), nil
}

// DeleteDocument deletes the document at the provvided id.
func (d *Database) DeleteDocument(ctx context.Context, id, rev string) error {
	request, err := d.client.requestFor(ctx, nil, http.MethodDelete, d.name, id)
	if err != nil {
		return err
	}
	request.Header.Set("If-Match", rev)

	response, err := d.client.do(request)
	if err != nil {
		return errx.Annotatef(err, "request [%s] [%s]", request.Method, request.URL)
	}
	defer response.Body.Close()

	_, err = evaluateResponse(response)
	if errx.IsNotFound(err) {
		return errx.NotFoundf("could not find document with id [missing] and revision []")
	}
	if err != nil {
		return err
	}

	return nil
}

// StoreAttachment stores the provided attachment at the provided document id.
func (d *Database) StoreAttachment(ctx context.Context, id, rev, name, contentType string, content []byte) (string, error) {
	request, err := d.client.requestFor(ctx, bytes.NewBuffer(content), http.MethodPut, d.name, id, name)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("If-Match", rev)
	request.Header.Set("Accept", MimeTypeJSON)

	response, err := d.client.do(request)
	if err != nil {
		return "", errx.Annotatef(err, "request [%s] [%s]", request.Method, request.URL)
	}
	defer response.Body.Close()

	m, err := evaluateResponse(response)
	if errx.IsAlreadyExists(err) {
		return "", errx.NotFoundf("could not find document with id [%s] and revision [%s]", id, rev)
	}
	if err != nil {
		return "", err
	}

	return m.(map[string]interface{})["rev"].(string), nil
}

// DeleteAttachment deletes the attachment for the provided document id.
func (d *Database) DeleteAttachment(ctx context.Context, id, rev, name string) error {
	request, err := d.client.requestFor(ctx, nil, http.MethodDelete, d.name, id, name)
	if err != nil {
		return err
	}
	request.Header.Set("If-Match", rev)
	request.Header.Set("X-Couch-Full-Commit", "true")

	response, err := d.client.do(request)
	if err != nil {
		return errx.Annotatef(err, "request [%s] [%s]", request.Method, request.URL)
	}
	defer response.Body.Close()

	_, err = evaluateResponse(response)
	if errx.IsAlreadyExists(err) {
		return errx.NotFoundf("could not find attachment with id [%s], revision [%s] and name [%s]", id, rev, name)
	}
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes the couch db database.
func (d *Database) Delete(ctx context.Context) error {
	request, err := d.client.requestFor(ctx, nil, http.MethodDelete, d.name)
	if err != nil {
		return err
	}

	response, err := d.client.do(request)
	if err != nil {
		return errx.Annotatef(err, "request [%s] [%s]", request.Method, request.URL)
	}
	defer response.Body.Close()

	if _, err := evaluateResponse(response); err != nil {
		return err
	}

	return nil
}
