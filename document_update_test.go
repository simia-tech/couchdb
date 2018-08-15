package couchdb_test

import (
	"testing"

	"github.com/simia-tech/couchdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocumentUpdate(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	db := e.client.Database("test")
	db.Create().Do()
	defer db.Delete().Do()

	document := couchdb.Document{"test": "value"}
	createResponse, err := db.Document().Create(document).Do()
	require.NoError(t, err)

	document["test"] = "another value"
	response, err := db.Document().Update(createResponse.ID, document).WithRevision(createResponse.Revision).Do()
	require.NoError(t, err)

	assert.True(t, response.OK)
	assert.Equal(t, createResponse.ID, response.ID)
	assert.NotEqual(t, createResponse.Revision, response.Revision)
}
