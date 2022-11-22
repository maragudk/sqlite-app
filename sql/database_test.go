package sql_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/maragudk/litefs-app/sqltest"
)

func TestDatabase_MigrateDown(t *testing.T) {
	t.Run("can migrate down", func(t *testing.T) {
		db := sqltest.CreateDatabase(t)
		err := db.MigrateDown(context.Background())
		require.NoError(t, err)
	})
}

func TestDatabase_GetPrimary(t *testing.T) {
	t.Run("returns the empty string when no .primary file", func(t *testing.T) {
		db := sqltest.CreateDatabase(t)

		primary, err := db.GetPrimary()
		require.NoError(t, err)
		require.Equal(t, "", primary)
	})

	t.Run("returns instance set in primary file", func(t *testing.T) {
		err := os.WriteFile(".primary", []byte("foo"), 0644)
		require.NoError(t, err)

		defer func() {
			err := os.Remove(".primary")
			require.NoError(t, err)
		}()

		db := sqltest.CreateDatabase(t)

		primary, err := db.GetPrimary()
		require.NoError(t, err)
		require.Equal(t, "foo", primary)
	})
}
