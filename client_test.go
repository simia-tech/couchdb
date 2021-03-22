package couchdb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientInfo(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	t.Run("InstanceInfo", func(t *testing.T) {
		ii, err := e.client.InstanceInfo(e.ctx)
		require.NoError(t, err)

		assert.Equal(t, "Welcome", ii.CouchDB)
		assert.Equal(t, "3.1.1", ii.Version)
		assert.Equal(t, "ce596c65d", ii.GitSHA)
		assert.Equal(t, []string{"access-ready", "partitioned", "pluggable-storage-engines", "reshard", "scheduler"},
			ii.Features)
		assert.Equal(t, "The Apache Software Foundation", ii.Vendor.Name)
	})

	t.Run("AllDatabases", func(t *testing.T) {
		dbs, err := e.client.AllDatabases(e.ctx)
		require.NoError(t, err)

		assert.Equal(t, []string{}, dbs)
	})
}
