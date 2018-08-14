package routes

import (
	"github.com/go-chi/chi"
	"github.com/gregbiv/sandbox/pkg/api/keys"
	"github.com/gregbiv/sandbox/pkg/storage"
	"github.com/jmoiron/sqlx"
)

// RouteKeys registers keys routes
func RouteKeys(db *sqlx.DB) func(r chi.Router) {
	getter := storage.NewGetter(db)

	return func(r chi.Router) {
		r.Get("/", keys.NewGetKeysHandler(getter).ServeHTTP)
	}
}
