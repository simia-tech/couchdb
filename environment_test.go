package couchdb_test

import (
	"context"
	"testing"

	"github.com/simia-tech/couchdb"
	"github.com/simia-tech/env"
)

var (
	couchDBUserURL  = env.String("COUCHDB_BASE_URL", "http://127.0.0.1:5984")
	couchDBUsername = env.String("COUCHDB_USERNAME", "admin")
	couchDBPassword = env.String("COUCHDB_PASSWORD", "")
)

type environment struct {
	ctx      context.Context
	client   *couchdb.Client
	tearDown func()
}

func setUpTestEnvironment(tb testing.TB) *environment {
	ctx := context.Background()

	client := &couchdb.Client{
		BaseURL:  couchDBUserURL.Get(),
		Username: couchDBUsername.Get(),
		Password: couchDBPassword.Get(),
	}

	return &environment{
		ctx:      ctx,
		client:   client,
		tearDown: func() {},
	}
}
