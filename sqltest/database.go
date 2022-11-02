package sqltest

import (
	"context"
	"os"
	"testing"

	"github.com/maragudk/env"
	"github.com/maragudk/migrate"

	"github.com/maragudk/sqlite-app/sql"
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
	if err := migrate.Up(context.Background(), db.DB.DB, os.DirFS("../sql/migrations")); err != nil {
		t.Fatal(err)
	}
	if err := migrate.Down(context.Background(), db.DB.DB, os.DirFS("../sql/migrations")); err != nil {
		t.Fatal(err)
	}
	if err := migrate.Up(context.Background(), db.DB.DB, os.DirFS("../sql/migrations")); err != nil {
		t.Fatal(err)
	}

	return db
}
