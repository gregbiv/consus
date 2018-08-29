package routes

import (
	"github.com/go-chi/chi"
	"github.com/gregbiv/consus/pkg/api/keys"
	"github.com/gregbiv/consus/pkg/storage"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
)

// RouteKeys registers keys routes
func RouteKeys(db *sqlx.DB) func(r chi.Router) {
	getter := storage.NewGetter(db)
	discarder := storage.NewDiscarder(db)
	storer := storage.NewStorer(db)
	updater := storage.NewUpdater(db)

	instrf := prometheus.InstrumentHandlerFunc

	return func(r chi.Router) {
		r.Get("/", instrf("get_keys", keys.NewGetKeysHandler(getter).ServeHTTP))
		r.Delete("/", instrf("delete_keys", keys.NewDiscardKeysHandler(discarder).ServeHTTP))
		r.Put("/", instrf("put_keys", keys.NewPutKeyHandler(storer, updater, getter).ServeHTTP))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", instrf("get_by_id", keys.NewGetKeyHandler(getter).ServeHTTP))
			r.Head("/", instrf("head_by_id", keys.NewHeadKeyHandler(getter).ServeHTTP))
			r.Delete("/", instrf("delete_by_id", keys.NewDiscardKeyHandler(discarder).ServeHTTP))
		})
	}
}
