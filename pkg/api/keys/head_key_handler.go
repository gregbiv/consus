package keys

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/gregbiv/sandbox/pkg/api"
	"github.com/gregbiv/sandbox/pkg/storage"
)

type (
	headKeyHandler struct {
		keyGetter storage.Getter
	}
)

// NewHeadKeyHandler init and returns an instance of headKeyHandler
func NewHeadKeyHandler(keyGetter storage.Getter) http.Handler {
	return &headKeyHandler{
		keyGetter: keyGetter,
	}
}

func (h *headKeyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "id")

	_, err := h.keyGetter.GetByID(ID, true)
	if err != nil {
		if err == storage.ErrKeyNotFound {
			api.NotFound(w, r)
			return
		}
		api.RenderInternalServerError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
}
