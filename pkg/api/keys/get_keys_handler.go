package keys

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/gregbiv/sandbox/pkg/api"
	"github.com/gregbiv/sandbox/pkg/storage"
)

type (
	getKeysHandler struct {
		keyGetter storage.Getter
	}

	ListResponse struct {
		Items []key `json:"items"`
	}
)

// NewGetKeysHandler init and returns an instance of getKeysHandler
func NewGetKeysHandler(KeyGetter storage.Getter) http.Handler {
	return &getKeysHandler{
		keyGetter: KeyGetter,
	}
}

func (h *getKeysHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// getting filter parameter
	filterStr := r.URL.Query().Get("filter")

	// fetching all keys
	dbKeys, err := h.keyGetter.GetAll(filterStr, true)
	if err != nil {
		api.RenderInternalServerError(w, r, err)
		return
	}

	list := []key{}
	for _, dbKey := range dbKeys {
		keyAPI := key{}
		keyAPI.fromDB(dbKey)

		list = append(list, keyAPI)
	}

	response := ListResponse{Items: list}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
