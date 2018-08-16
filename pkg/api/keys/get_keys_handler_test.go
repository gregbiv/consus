package keys

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gregbiv/sandbox/pkg/api"
	"github.com/gregbiv/sandbox/pkg/model"
	"github.com/gregbiv/sandbox/pkg/storage/mock"
	"github.com/stretchr/testify/assert"
)

func TestGetKeysHandler_ServeHTTP(t *testing.T) {
	t.Parallel()

	// creating fixtures
	keyModel := model.KeyFactory.MustCreate().(*model.Key)
	// creating storage mocks
	keyGetter := &mock.Getter{Key: keyModel}
	// creating handler
	handler := NewGetKeysHandler(keyGetter)

	t.Run("Getting all keys", func(t *testing.T) {
		var keyListResponse listResponse

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/keys/", nil)

		handler.ServeHTTP(w, r)
		err := json.Unmarshal(w.Body.Bytes(), &keyListResponse)

		assert.Nil(t, err)
		assert.IsType(t, listResponse{}, keyListResponse)
		assert.Equal(t, 1, len(keyListResponse.Items))
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetKeysHandler_ServeHTTP_Negative(t *testing.T) {
	// creating storage mocks
	keyGetter := &mock.Getter{}
	// creating handler
	handler := NewGetKeysHandler(keyGetter)

	t.Run("Getting a key without providing an ID", func(t *testing.T) {
		var keyListResponse listResponse

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/keys", nil)

		handler.ServeHTTP(w, r)
		err := json.Unmarshal(w.Body.Bytes(), &keyListResponse)

		assert.Nil(t, err)
		assert.IsType(t, listResponse{}, keyListResponse)
		assert.Equal(t, 0, len(keyListResponse.Items))
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetKeysHandler_ServeHTTP_Failure(t *testing.T) {
	// creating storage mocks
	keyGetter := &mock.Getter{KeyError: errors.New("failed to fetch data from DB")}
	// creating handler
	handler := NewGetKeysHandler(keyGetter)

	t.Run("Getting a key ", func(t *testing.T) {
		var errResponse api.ErrResponse

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/keys", nil)

		handler.ServeHTTP(w, r)
		err := json.Unmarshal(w.Body.Bytes(), &errResponse)

		assert.Nil(t, err)
		assert.IsType(t, api.ErrResponse{}, errResponse)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
