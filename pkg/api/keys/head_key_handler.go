package keys

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/gregbiv/sandbox/pkg/api"
	"github.com/gregbiv/sandbox/pkg/storage"
)

type (
	headKeyHandler struct {
		urlExtractor api.URLExtractor
		keyGetter    storage.Getter
	}
)

// NewHeadKeyHandler init and returns an instance of headKeyHandler
func NewHeadKeyHandler(urlExtractor api.URLExtractor, keyGetter storage.Getter) http.Handler {
	return &headKeyHandler{
		urlExtractor: urlExtractor,
		keyGetter:    keyGetter,
	}
}

func (h *headKeyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ID, err := h.urlExtractor.UUIDFromRoute(r, "id")
	if err != nil {
		api.NotFound(w, r)
		return
	}

	_, err = h.keyGetter.GetByID(ID.String())
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
