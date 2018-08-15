package keys

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/render"
	"github.com/gregbiv/sandbox/pkg/api"
	"github.com/gregbiv/sandbox/pkg/storage"
)

type (
	putKeyHandler struct {
		keyStorer  storage.Storer
		keyUpdater storage.Updater
		keyGetter  storage.Getter
	}
)

// NewPutKeyHandler init and returns an instance of putKeyHandler
func NewPutKeyHandler(keyStorer storage.Storer, keyUpdater storage.Updater, keyGetter storage.Getter) http.Handler {
	return &putKeyHandler{
		keyStorer:  keyStorer,
		keyUpdater: keyUpdater,
		keyGetter:  keyGetter,
	}
}

func (h *putKeyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// This key api model will hold all
	// the params from the request
	keyAPI := key{}

	// getting expire_in parameter
	expireIn := r.URL.Query().Get("expire_in")
	if expireIn != "" {
		seconds, err := strconv.Atoi(expireIn)

		if err != nil || seconds <= 0 {
			api.RenderInvalidInput(w, r, "expire_in", err.Error())
			return
		}

		timeIn := time.Now().Add(time.Duration(seconds) * time.Second)
		keyAPI.ExpiresAt = &timeIn
	}

	err := keyAPI.fromRequest(r)
	if err != nil {
		if err == ErrInvalidBody {
			api.RenderInvalidInput(w, r, "", ErrInvalidBody.Error())
			return
		}
		api.RenderInternalServerError(w, r, err)
		return
	}

	modelKey, err := keyAPI.toModel()
	if err != nil {
		api.RenderInvalidInput(w, r, "value", err.Error())
		return
	}

	dbKey, err := h.keyGetter.GetByID(modelKey.KeyID)
	if err != nil {
		if err != storage.ErrKeyNotFound {
			api.RenderInternalServerError(w, r, err)
			return
		}

		if err := h.keyStorer.Store(modelKey); err != nil {
			api.RenderInternalServerError(w, r, err)
			return
		}
	} else {
		dbKey.Value = modelKey.Value
		dbKey.ExpiresAt = modelKey.ExpiresAt

		if err := h.keyUpdater.Update(dbKey); err != nil {
			api.RenderInternalServerError(w, r, err)
			return
		}
	}

	render.Status(r, http.StatusOK)
}
