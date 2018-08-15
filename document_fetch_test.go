package couchdb_test

import (
	"testing"

	"github.com/simia-tech/couchdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocumentFetch(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	db := e.client.Database("test")
	db.Create().Do()
	defer db.Delete().Do()

	document := couchdb.Document{"test": "value"}
	createResponse, err := db.Document().Create(document).Do()
	require.NoError(t, err)
	document["_id"] = createResponse.ID
	document["_rev"] = createResponse.Revision

	response, err := db.Document().Fetch(createResponse.ID).WithRevision(createResponse.Revision).Do()
	require.NoError(t, err)

	assert.Equal(t, createResponse.Revision, response.Revision)
	assert.Equal(t, document, response.Document)
}
