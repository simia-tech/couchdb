package couchdb_test

import (
	"testing"

	"github.com/simia-tech/couchdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocumentFetchMeta(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	db := e.client.Database("test")
	db.Create().Do()
	defer db.Delete().Do()

	document := couchdb.Document{"test": "value"}
	createResponse, err := db.Document().Create(document).Do()
	require.NoError(t, err)

	response, err := db.Document().FetchMeta(createResponse.ID).WithRevision(createResponse.Revision).Do()
	require.NoError(t, err)

	assert.Equal(t, uint64(0x66), response.ContentLength)
	assert.Equal(t, createResponse.Revision, response.Revision)
}
