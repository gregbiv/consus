package routes

import (
	"testing"

	"github.com/go-chi/chi"
	"github.com/gregbiv/consus/pkg/storage/mock"
	"github.com/stretchr/testify/assert"
)

func TestRouteKeys(t *testing.T) {
	db, _, _ := mock.NewDbMock()
	defer db.Close()

	router := chi.NewRouter()
	router.Route("/keys", RouteKeys(db))

	assert.Equal(t, 1, len(router.Routes()))
}
