package couchdb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseCreate(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()
	db := e.client.Database("test")

	response, err := db.Create().Do()
	require.NoError(t, err)
	defer db.Delete().Do()

	assert.True(t, response.OK)
	assert.Equal(t, "", response.Error)
	assert.Equal(t, "", response.Reason)
}