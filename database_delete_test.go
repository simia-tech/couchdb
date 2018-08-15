package couchdb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseDelete(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	db := e.client.Database("test")
	_, err := db.Create().Do()
	require.NoError(t, err)

	response, err := db.Delete().Do()
	require.NoError(t, err)

	assert.True(t, response.OK)
}
