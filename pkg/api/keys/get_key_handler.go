package keys

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/gregbiv/sandbox/pkg/api"
	"github.com/gregbiv/sandbox/pkg/storage"
)

type (
	getKeyHandler struct {
		keyGetter storage.Getter
	}
)

// NewGetKeyHandler init and returns an instance of getKeyHandler
func NewGetKeyHandler(keyGetter storage.Getter) http.Handler {
	return &getKeyHandler{
		keyGetter: keyGetter,
	}
}

func (h *getKeyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "id")

	dbKey, err := h.keyGetter.GetByID(ID, true)
	if err != nil {
		if err == storage.ErrKeyNotFound {
			api.NotFound(w, r)
			return
		}
		api.RenderInternalServerError(w, r, err)
		return
	}

	keyAPI := key{}
	keyAPI.fromDB(dbKey)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, keyAPI)
}
