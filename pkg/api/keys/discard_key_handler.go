package keys

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/gregbiv/sandbox/pkg/api"
	"github.com/gregbiv/sandbox/pkg/storage"
)

type (
	discardKeyHandler struct {
		urlExtractor api.URLExtractor
		keyDiscarder storage.Discarder
	}
)

// NewDiscardKeyHandler init and returns an instance of discardKeyHandler
func NewDiscardKeyHandler(urlExtractor api.URLExtractor, keyDiscarder storage.Discarder) http.Handler {
	return &discardKeyHandler{
		urlExtractor: urlExtractor,
		keyDiscarder: keyDiscarder,
	}
}

func (h *discardKeyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ID, err := h.urlExtractor.UUIDFromRoute(r, "id")
	if err != nil {
		api.NotFound(w, r)
		return
	}

	err = h.keyDiscarder.Discard(ID.String())
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
