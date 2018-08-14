package couchdb_test

import (
	"testing"

	"code.posteo.de/common/errx"
	"github.com/stretchr/testify/require"
)

func TestClientCreateAndDeleteDatabase(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	db, err := e.client.CreateDatabase(e.ctx, "test")
	require.NoError(t, err)

	err = db.Delete(e.ctx)
	require.NoError(t, err)
}

func TestClientCreateExistingDatabase(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	db, err := e.client.CreateDatabase(e.ctx, "test")
	require.NoError(t, err)

	_, err = e.client.CreateDatabase(e.ctx, "test")
	errx.AssertError(t, errx.AlreadyExistsf("database [test] already exists"), err)

	err = db.Delete(e.ctx)
	require.NoError(t, err)
}
