package couchdb_test

import (
	"testing"

	"code.posteo.de/common/errx"
	"github.com/simia-tech/couchdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseCreateDocument(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	db := e.client.GetDatabase("test")
	if db != nil {
		db.Delete(e.ctx)
	}

	db, err := e.client.CreateDatabase(e.ctx, "test")
	require.NoError(t, err)
	defer db.Delete(e.ctx)

	tcs := []struct {
		name             string
		document         map[string]interface{}
		expectError      error
		expectIDPattern  string
		expectRevPattern string
	}{
		{"Success", map[string]interface{}{"test": "value"}, nil, "[0-9a-f]{32}", "1-[0-9a-f]{32}"},
		{"MissingDocument", nil, errx.BadRequestf("Document must be a JSON object"), "", ""},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			id, rev, err := db.CreateDocument(e.ctx, tc.document)
			errx.AssertError(t, tc.expectError, err)
			assert.Regexp(t, tc.expectIDPattern, id)
			assert.Regexp(t, tc.expectRevPattern, rev)
		})
	}
}

func TestDatabaseFetchDocument(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	db, err := e.client.CreateDatabase(e.ctx, "test")
	require.NoError(t, err)
	defer db.Delete(e.ctx)

	document := couchdb.Document{"test": "value"}
	id, rev, err := db.CreateDocument(e.ctx, document)
	require.NoError(t, err)
	document["_id"] = id
	document["_rev"] = rev

	tcs := []struct {
		name           string
		id             string
		rev            string
		expectError    error
		expectDocument couchdb.Document
	}{
		{"Success", id, rev, nil, document},
		{"Latest", id, "", nil, document},
		{"MissingDocument", "missing", "", errx.NotFoundf("could not find document with id [missing] and revision []"), nil},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			document, err := db.FetchDocument(e.ctx, tc.id, tc.rev)
			errx.AssertError(t, tc.expectError, err)

			assert.Equal(t, tc.expectDocument, document)
		})
	}
}

func TestDatabaseFetchDocumentMeta(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	db, err := e.client.CreateDatabase(e.ctx, "test")
	require.NoError(t, err)
	defer db.Delete(e.ctx)

	document := map[string]interface{}{"test": "value"}
	id, rev, err := db.CreateDocument(e.ctx, document)
	require.NoError(t, err)
	document["_id"] = id

	tcs := []struct {
		name                string
		id                  string
		expectError         error
		expectRevision      string
		expectContentLength uint64
	}{
		{"Success", id, nil, rev, 0x66},
		{"MissingDocument", "missing", errx.NotFoundf("could not find document revision with id [missing]"), "", 0},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			revision, contentLength, err := db.FetchDocumentMeta(e.ctx, tc.id)
			errx.AssertError(t, tc.expectError, err)

			assert.Equal(t, tc.expectRevision, revision)
			assert.Equal(t, tc.expectContentLength, contentLength)
		})
	}
}

func TestDatabaseUpdateDocument(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	db, err := e.client.CreateDatabase(e.ctx, "test")
	require.NoError(t, err)
	defer db.Delete(e.ctx)

	document := map[string]interface{}{"test": "value"}
	id, rev, err := db.CreateDocument(e.ctx, document)
	require.NoError(t, err)

	newDocument := map[string]interface{}{"test": "another value"}

	tcs := []struct {
		name                string
		id                  string
		rev                 string
		document            map[string]interface{}
		expectError         error
		expectDocumentValue string
	}{
		{"Success", id, rev, newDocument, nil, "another value"},
		{"MissingDocument", "missing", rev, newDocument, nil, "another value"},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			newRev, err := db.UpdateDocument(e.ctx, tc.id, tc.rev, tc.document)
			errx.AssertError(t, tc.expectError, err)

			if tc.expectError == nil {
				document, err := db.FetchDocument(e.ctx, tc.id, newRev)
				require.NoError(t, err)

				assert.Equal(t, tc.expectDocumentValue, document["test"])
			}
		})
	}
}

func TestDatabaseDeleteDocument(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	db, err := e.client.CreateDatabase(e.ctx, "test")
	require.NoError(t, err)
	defer db.Delete(e.ctx)

	document := map[string]interface{}{"test": "value"}
	id, rev, err := db.CreateDocument(e.ctx, document)
	require.NoError(t, err)

	tcs := []struct {
		name        string
		id          string
		rev         string
		expectError error
	}{
		{"Success", id, rev, nil},
		{"MissingDocument", "missing", rev, errx.NotFoundf("could not find document with id [missing] and revision []")},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := db.DeleteDocument(e.ctx, tc.id, tc.rev)
			errx.AssertError(t, tc.expectError, err)
		})
	}
}

func TestDatabaseStoreAttachment(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	db, err := e.client.CreateDatabase(e.ctx, "test")
	require.NoError(t, err)
	defer db.Delete(e.ctx)

	id, rev, err := db.CreateDocument(e.ctx, map[string]interface{}{"test": "value"})
	require.NoError(t, err)

	tcs := []struct {
		name                  string
		id                    string
		rev                   string
		attachmentName        string
		attachmentContentType string
		attachmentContent     []byte
		expectError           error
		expectNewRev          bool
	}{
		{"Success", id, rev, "test", "application/octet-stream", []byte{0, 1, 2, 3},
			nil, true},
		{"MissingRevision", id, "", "test", "application/octet-stream", []byte{0, 1, 2, 3},
			errx.BadRequestf("Invalid rev format"), false},
		{"InvalidRevision", id, "1-5b497a0860ecee16f567f88e7ede1883", "test", "application/octet-stream", []byte{0, 1, 2, 3},
			errx.NotFoundf("could not find document with id [%s] and revision [1-5b497a0860ecee16f567f88e7ede1883]", id), false},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			newRev, err := db.StoreAttachment(e.ctx, tc.id, tc.rev, tc.attachmentName, tc.attachmentContentType, tc.attachmentContent)
			errx.AssertError(t, tc.expectError, err)
			if tc.expectNewRev {
				assert.NotEqual(t, tc.rev, newRev)
			}
		})
	}
}

func TestDatabaseDeleteAttachment(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	db, err := e.client.CreateDatabase(e.ctx, "test")
	require.NoError(t, err)
	defer db.Delete(e.ctx)

	id, rev, err := db.CreateDocument(e.ctx, map[string]interface{}{"test": "value"})
	require.NoError(t, err)

	tcs := []struct {
		name           string
		attachmentName string
		expectError    error
	}{
		{"Success", "test", nil},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			rev, err = db.StoreAttachment(e.ctx, id, rev, "test", "application/octet-stream", []byte{0, 1, 2, 3})
			require.NoError(t, err)

			err = db.DeleteAttachment(e.ctx, id, rev, tc.attachmentName)
			errx.AssertError(t, tc.expectError, err)
		})
	}
}
