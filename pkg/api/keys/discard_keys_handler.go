package keys

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/gregbiv/sandbox/pkg/api"
	"github.com/gregbiv/sandbox/pkg/storage"
)

type (
	discardKeysHandler struct {
		keyDiscarder storage.Discarder
	}
)

// NewDiscardKeyHandler init and returns an instance of discardKeyHandler
func NewDiscardKeysHandler(keyDiscarder storage.Discarder) http.Handler {
	return &discardKeysHandler{
		keyDiscarder: keyDiscarder,
	}
}

func (h *discardKeysHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.keyDiscarder.Truncate()
	if err != nil {
		api.RenderInternalServerError(w, r, err)
		return
	}

	render.NoContent(w, r)
}
