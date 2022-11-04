package sql_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/maragudk/sqlite-app/sqltest"
)

func TestDatabase_MigrateDown(t *testing.T) {
	t.Run("can migrate down", func(t *testing.T) {
		db := sqltest.CreateDatabase(t)
		err := db.MigrateDown(context.Background())
		require.NoError(t, err)
	})
}
