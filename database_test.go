package couchdb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/couchdb"
)

func TestDatabase(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	t.Run("Create", func(t *testing.T) {
		db := couchdb.NewDatabase(e.client, "test")
		require.NoError(t, db.Create(e.ctx))
		defer db.Delete(e.ctx)
	})

	t.Run("CreateExisting", func(t *testing.T) {
		db := couchdb.NewDatabase(e.client, "test")
		require.NoError(t, db.Create(e.ctx))
		defer db.Delete(e.ctx)

		err := db.Create(e.ctx)
		assert.ErrorIs(t, err, couchdb.ErrDatabaseAlreadyExists)
	})

	t.Run("Delete", func(t *testing.T) {
		db := couchdb.NewDatabase(e.client, "test")
		require.NoError(t, db.Create(e.ctx))

		require.NoError(t, db.Delete(e.ctx))
	})

	t.Run("DeleteMissing", func(t *testing.T) {
		db := couchdb.NewDatabase(e.client, "test")

		err := db.Delete(e.ctx)
		assert.ErrorIs(t, err, couchdb.ErrDatabaseDoesNotExists)
	})

	t.Run("Info", func(t *testing.T) {
		db := couchdb.NewDatabase(e.client, "test")
		require.NoError(t, db.Create(e.ctx))
		defer db.Delete(e.ctx)

		di, err := db.Info(e.ctx)
		require.NoError(t, err)

		assert.Equal(t, "test", di.Name)
		assert.Equal(t, uint(1), di.Cluster.N)
		assert.Equal(t, uint(2), di.Cluster.Q)
		assert.Equal(t, uint(1), di.Cluster.W)
		assert.Equal(t, uint(1), di.Cluster.R)
	})

}
