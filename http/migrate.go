package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type migrator interface {
	MigrateUp(ctx context.Context) error
}

func Migrate(mux chi.Router, db migrator) {
	mux.Post("/migrate/up", func(w http.ResponseWriter, r *http.Request) {
		if err := db.MigrateUp(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
