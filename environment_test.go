package couchdb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/simia-tech/couchdb"
)

type environment struct {
	ctx      context.Context
	client   *couchdb.Client
	tearDown func()
}

func setUpTestEnvironment(tb testing.TB) *environment {
	ctx := context.Background()

	client, err := couchdb.NewClient("http://127.0.0.1:5984", couchdb.WithUsername("admin"), couchdb.WithPassword("admin"))
	require.NoError(tb, err)

	return &environment{
		ctx:      ctx,
		client:   client,
		tearDown: func() {},
	}
}
