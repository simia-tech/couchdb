package couchdb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/couchdb"
)

func TestDocument(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	db := couchdb.NewDatabase(e.client, "test")
	require.NoError(t, db.Create(e.ctx))
	defer db.Delete(e.ctx)

	t.Run("Store", func(t *testing.T) {
		t.Run("WithoutID", func(t *testing.T) {
			document := couchdb.NewDocument(db, "", "")

			require.NoError(t, document.Store(e.ctx, map[string]interface{}{"test": "value"}))

			assert.Regexp(t, `^[0-9a-f]+$`, document.ID())
			assert.Regexp(t, `^\d+\-[0-9a-f]+$`, document.Revision())
		})

		t.Run("WithID", func(t *testing.T) {
			document := couchdb.NewDocument(db, "test", "")

			require.NoError(t, document.Store(e.ctx, map[string]interface{}{"test": "value"}))

			assert.Equal(t, "test", document.ID())
			assert.Regexp(t, `^\d+\-[0-9a-f]+$`, document.Revision())
		})

		t.Run("WithIDAndRevision", func(t *testing.T) {
			document := couchdb.NewDocument(db, "", "")
			require.NoError(t, document.Store(e.ctx, map[string]interface{}{"test": "value"}))

			id, revision := document.ID(), document.Revision()

			require.NoError(t, document.Store(e.ctx, map[string]interface{}{"test": "another value"}))

			assert.Equal(t, id, document.ID())
			assert.NotEqual(t, revision, document.Revision())
		})

		t.Run("WithIDAndMalformedRevision", func(t *testing.T) {
			document := couchdb.NewDocument(db, "", "")
			require.NoError(t, document.Store(e.ctx, map[string]interface{}{"test": "value"}))

			otherDocument := couchdb.NewDocument(db, document.ID(), "malformed")

			err := otherDocument.Store(e.ctx, map[string]interface{}{"test": "another value"})
			assert.ErrorIs(t, err, couchdb.ErrBadRequest)
		})

		t.Run("WithIDAndInvalidRevision", func(t *testing.T) {
			document := couchdb.NewDocument(db, "", "")
			require.NoError(t, document.Store(e.ctx, map[string]interface{}{"test": "value"}))

			otherDocument := couchdb.NewDocument(db, document.ID(), "0-123")

			err := otherDocument.Store(e.ctx, map[string]interface{}{"test": "another value"})
			assert.ErrorIs(t, err, couchdb.ErrConflict)
		})
	})

	t.Run("Fetch", func(t *testing.T) {
		d := couchdb.NewDocument(db, "", "")
		require.NoError(t, d.Store(e.ctx, map[string]interface{}{"test": "value"}))

		t.Run("WithID", func(t *testing.T) {
			document := couchdb.NewDocument(db, d.ID(), "")

			data := map[string]interface{}{}
			require.NoError(t, document.Fetch(e.ctx, &data))

			assert.Equal(t, d.ID(), document.ID())
			assert.Regexp(t, `^\d+\-[0-9a-f]+$`, document.Revision())

			assert.Equal(t, map[string]interface{}{
				"_id":  d.ID(),
				"_rev": d.Revision(),
				"test": "value",
			}, data)
		})
	})
}
