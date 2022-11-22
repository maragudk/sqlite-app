package sqltest

import (
	"context"
	"testing"

	"github.com/maragudk/env"

	"github.com/maragudk/litefs-app/sql"
)

// CreateDatabase for testing.
func CreateDatabase(t *testing.T) *sql.Database {
	t.Helper()

	_ = env.Load("../.env-test")

	db := sql.NewDatabase(sql.NewDatabaseOptions{
		URL:                env.GetStringOrDefault("DATABASE_URL", ":memory:"),
		MaxOpenConnections: 1,
		MaxIdleConnections: 1,
	})
	if err := db.Connect(); err != nil {
		t.Fatal(err)
	}

	if err := db.MigrateUp(context.Background()); err != nil {
		t.Fatal(err)
	}

	return db
}
