package keys

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/gregbiv/sandbox/pkg/api"
	"github.com/gregbiv/sandbox/pkg/storage"
)

type (
	getKeyHandler struct {
		urlExtractor api.URLExtractor
		keyGetter    storage.Getter
	}
)

// NewGetKeyHandler init and returns an instance of getKeyHandler
func NewGetKeyHandler(urlExtractor api.URLExtractor, keyGetter storage.Getter) http.Handler {
	return &getKeyHandler{
		urlExtractor: urlExtractor,
		keyGetter:    keyGetter,
	}
}

func (h *getKeyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ID, err := h.urlExtractor.UUIDFromRoute(r, "id")
	if err != nil {
		api.NotFound(w, r)
		return
	}

	dbKey, err := h.keyGetter.GetByID(ID.String())
	if err != nil {
		if err == storage.ErrKeyNotFound {
			api.NotFound(w, r)
			return
		}
		api.RenderInternalServerError(w, r, err)
		return
	}

	keyAPI := key{}
	err = keyAPI.fromDB(dbKey)
	if err != nil {
		api.RenderInternalServerError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, keyAPI)
}
