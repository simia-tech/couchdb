package couchdb

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/simia-tech/couchdb/value"
)

// Various errors.
var (
	ErrDatabaseAlreadyExists = errors.New("database already exists")
	ErrDatabaseDoesNotExists = errors.New("database does not exists")
)

// Database implements all methods on a couchdb database.
type Database struct {
	client *Client
	name   string
}

// NewDatabase returns a new database using the provided `*Client` with the provided name.
func NewDatabase(client *Client, name string) *Database {
	return &Database{
		client: client,
		name:   name,
	}
}

// Create creates the database.
func (db *Database) Create(ctx context.Context) error {
	r := value.Status{}
	if err := db.client.requestJSON(ctx, http.MethodPut, "/"+db.name, &r); err != nil {
		return err
	}
	if !r.OK {
		switch r.Error {
		case "file_exists":
			return fmt.Errorf("create database %s: %w", db.name, ErrDatabaseAlreadyExists)
		default:
			return fmt.Errorf("unknown error %s: %s", r.Error, r.Reason)
		}
	}
	return nil
}

// Delete deletes the database.
func (db *Database) Delete(ctx context.Context) error {
	r := value.Status{}
	if err := db.client.requestJSON(ctx, http.MethodDelete, "/"+db.name, &r); err != nil {
		return err
	}
	if !r.OK {
		switch r.Error {
		case "not_found":
			return fmt.Errorf("delete database %s: %w", db.name, ErrDatabaseDoesNotExists)
		default:
			return fmt.Errorf("unknown error %s: %s", r.Error, r.Reason)
		}
	}
	return nil
}

// Info fetches infos about the database.
func (db *Database) Info(ctx context.Context) (value.DatabaseInfo, error) {
	r := value.DatabaseInfo{}
	if err := db.client.requestJSON(ctx, http.MethodGet, "/"+db.name, &r); err != nil {
		return r, err
	}
	return r, nil
}
