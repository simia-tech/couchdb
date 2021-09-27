package couchdb

import (
	"context"
	"errors"
	"net/http"
)

// Various errors.
var (
	ErrMissingID = errors.New("missing id")
)

// Document implements all methods on a couchdb document.
type Document struct {
	database *Database
	id       string
	revision string
}

// NewDocument returns a new document.
func NewDocument(database *Database, id, revision string) *Document {
	return &Document{
		database: database,
		id:       id,
		revision: revision,
	}
}

// ID returns the document's id.
func (d *Document) ID() string {
	return d.id
}

// Revision returns the document's revision.
func (d *Document) Revision() string {
	return d.revision
}

// Store saves the document to the database. If the id is empty, it will be save at
// a generated id.
func (d *Document) Store(ctx context.Context, data interface{}) error {
	if d.id == "" {
		return d.storeWithoutID(ctx, data)
	}
	return d.storeWithID(ctx, data)
}

// Fetch loads the document from the database.
func (d *Document) Fetch(ctx context.Context, data interface{}) error {
	if d.id == "" {
		return ErrMissingID
	}

	if err := d.database.client.requestJSON(ctx, http.MethodGet, "/"+d.database.name+"/"+d.id, nil, nil, data); err != nil {
		return err
	}

	return nil
}

func (d *Document) storeWithoutID(ctx context.Context, data interface{}) error {
	r := struct {
		OK       bool   `json:"ok"`
		ID       string `json:"id"`
		Revision string `json:"rev"`
	}{}
	if err := d.database.client.requestJSON(ctx, http.MethodPost, "/"+d.database.name, nil, data, &r); err != nil {
		return err
	}
	if r.OK {
		d.id = r.ID
		d.revision = r.Revision
	}
	return nil
}

func (d *Document) storeWithID(ctx context.Context, data interface{}) error {
	header := http.Header{}
	if d.revision != "" {
		header.Add("If-Match", d.revision)
	}
	r := struct {
		OK       bool   `json:"ok"`
		ID       string `json:"id"`
		Revision string `json:"rev"`
	}{}
	if err := d.database.client.requestJSON(ctx, http.MethodPut, "/"+d.database.name+"/"+d.id, header, data, &r); err != nil {
		return err
	}
	if r.OK {
		d.revision = r.Revision
	}
	return nil
}
