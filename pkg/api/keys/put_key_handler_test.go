package keys

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gregbiv/consus/pkg/model"

	"github.com/gregbiv/consus/pkg/api"

	"github.com/gregbiv/consus/pkg/storage"

	"github.com/gregbiv/consus/pkg/storage/mock"
	"github.com/stretchr/testify/assert"
)

func TestPutKeyHandler_ServeHTTP(t *testing.T) {
	t.Parallel()

	// creating storage mocks
	keyGetter := &mock.Getter{KeyError: storage.ErrKeyNotFound}
	keyUpdater := &mock.Updater{}
	keyStorer := &mock.Storer{}
	// creating handler
	handler := NewPutKeyHandler(keyStorer, keyUpdater, keyGetter)

	t.Run("Creating a new key by providing the payload", func(t *testing.T) {
		payload := `
		{
			"id": "test",
			"value": "some random text"
		}
		`
		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/keys", strings.NewReader(payload))

		handler.ServeHTTP(w, r)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Creating a new key by providing the payload and expire_in", func(t *testing.T) {
		payload := `
		{
			"id": "test",
			"value": "some random text"
		}
		`
		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/keys?expire_in=60", strings.NewReader(payload))

		handler.ServeHTTP(w, r)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Creating a new key by providing the payload and invalid expire_in", func(t *testing.T) {
		var errResponse api.ErrResponse
		payload := `
		{
			"id": "test",
			"value": "some random text"
		}
		`
		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/keys?expire_in=invalid", strings.NewReader(payload))

		handler.ServeHTTP(w, r)
		err := json.Unmarshal(w.Body.Bytes(), &errResponse)

		assert.Nil(t, err)
		assert.IsType(t, api.ErrResponse{}, errResponse)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "expire_in", errResponse.Errors.Target)
	})

	t.Run("Creating a new key by providing an incorrect payload", func(t *testing.T) {
		var errResponse api.ErrResponse
		payload := `
		{
			"id": "this is not a valid key",
			"value": "some random text"
		}
		`
		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/keys", strings.NewReader(payload))

		handler.ServeHTTP(w, r)
		err := json.Unmarshal(w.Body.Bytes(), &errResponse)

		assert.Nil(t, err)
		assert.IsType(t, api.ErrResponse{}, errResponse)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "id", errResponse.Errors.Target)
	})

	t.Run("Creating a new key by providing an empty payload", func(t *testing.T) {
		var errResponse api.ErrResponse
		payload := `{}`

		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/keys", strings.NewReader(payload))

		handler.ServeHTTP(w, r)
		err := json.Unmarshal(w.Body.Bytes(), &errResponse)

		assert.Nil(t, err)
		assert.IsType(t, api.ErrResponse{}, errResponse)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "id", errResponse.Errors.Target)
	})

	t.Run("Creating a new key without providing a payload", func(t *testing.T) {
		var errResponse api.ErrResponse

		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/keys", nil)

		handler.ServeHTTP(w, r)
		err := json.Unmarshal(w.Body.Bytes(), &errResponse)

		assert.Nil(t, err)
		assert.IsType(t, api.ErrResponse{}, errResponse)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Failing to read body while creating a new key", func(t *testing.T) {
		var errResponse api.ErrResponse
		// simulating invalid body
		file, _ := os.Open("my_file.zip")

		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/keys", file)

		handler.ServeHTTP(w, r)
		err := json.Unmarshal(w.Body.Bytes(), &errResponse)

		assert.Nil(t, err)
		assert.IsType(t, api.ErrResponse{}, errResponse)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestPutKeyHandler_ServeHTTP_GetterFailure(t *testing.T) {
	t.Parallel()

	// creating storage mocks
	keyGetter := &mock.Getter{KeyError: errors.New("failed to fetch data from DB")}
	keyUpdater := &mock.Updater{}
	keyStorer := &mock.Storer{}
	// creating handler
	handler := NewPutKeyHandler(keyStorer, keyUpdater, keyGetter)

	t.Run("Creating a new key by providing the payload", func(t *testing.T) {
		var errResponse api.ErrResponse
		payload := `
		{
			"id": "test",
			"value": "some random text"
		}
		`
		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/keys", strings.NewReader(payload))

		handler.ServeHTTP(w, r)
		err := json.Unmarshal(w.Body.Bytes(), &errResponse)

		assert.Nil(t, err)
		assert.IsType(t, api.ErrResponse{}, errResponse)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestPutKeyHandler_ServeHTTP_StoreFailure(t *testing.T) {
	t.Parallel()

	// creating storage mocks
	keyGetter := &mock.Getter{KeyError: storage.ErrKeyNotFound}
	keyUpdater := &mock.Updater{}
	keyStorer := &mock.Storer{StoreError: errors.New("failed to fetch data from DB")}
	// creating handler
	handler := NewPutKeyHandler(keyStorer, keyUpdater, keyGetter)

	t.Run("Creating a new key by providing the payload", func(t *testing.T) {
		var errResponse api.ErrResponse
		payload := `
		{
			"id": "test",
			"value": "some random text"
		}
		`
		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/keys", strings.NewReader(payload))

		handler.ServeHTTP(w, r)
		err := json.Unmarshal(w.Body.Bytes(), &errResponse)

		assert.Nil(t, err)
		assert.IsType(t, api.ErrResponse{}, errResponse)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestPutKeyHandler_ServeHTTP_UpdateFailure(t *testing.T) {
	t.Parallel()

	// creating fixtures
	keyModel := model.KeyFactory.MustCreate().(*model.Key)
	// creating storage mocks
	keyGetter := &mock.Getter{Key: keyModel}
	keyUpdater := &mock.Updater{UpdateError: errors.New("failed to fetch data from DB")}
	keyStorer := &mock.Storer{}
	// creating handler
	handler := NewPutKeyHandler(keyStorer, keyUpdater, keyGetter)

	t.Run("Creating a new key by providing the payload", func(t *testing.T) {
		var errResponse api.ErrResponse
		payload := `
		{
			"id": "test",
			"value": "some random text"
		}
		`
		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/keys", strings.NewReader(payload))

		handler.ServeHTTP(w, r)
		err := json.Unmarshal(w.Body.Bytes(), &errResponse)

		assert.Nil(t, err)
		assert.IsType(t, api.ErrResponse{}, errResponse)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
