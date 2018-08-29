package keys

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/gregbiv/consus/pkg/api"
	"github.com/gregbiv/consus/pkg/storage"
)

type (
	discardKeyHandler struct {
		keyDiscarder storage.Discarder
	}
)

// NewDiscardKeyHandler init and returns an instance of discardKeyHandler
func NewDiscardKeyHandler(keyDiscarder storage.Discarder) http.Handler {
	return &discardKeyHandler{
		keyDiscarder: keyDiscarder,
	}
}

func (h *discardKeyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "id")

	err := h.keyDiscarder.Discard(ID)
	if err != nil {
		if err == storage.ErrKeyNotFound {
			api.NotFound(w, r)
			return
		}
		api.RenderInternalServerError(w, r, err)
		return
	}

	render.NoContent(w, r)
}
