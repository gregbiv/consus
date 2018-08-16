package keys

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/gregbiv/sandbox/pkg/api"
	"github.com/gregbiv/sandbox/pkg/model"
	"github.com/gregbiv/sandbox/pkg/storage"
	"github.com/gregbiv/sandbox/pkg/storage/mock"
	"github.com/stretchr/testify/assert"
)

func TestGetKeyHandler_ServeHTTP(t *testing.T) {
	t.Parallel()

	// creating fixtures
	keyModel := model.KeyFactory.MustCreate().(*model.Key)
	// creating storage mocks
	keyGetter := &mock.Getter{Key: keyModel}
	// creating handler
	handler := NewGetKeyHandler(keyGetter)
	// Populate the request's context with our test data.
	rctx := chi.NewRouteContext()

	t.Run("Getting a key", func(t *testing.T) {
		var keyResponse key
		testKeyID := "test-bla-bla"

		rctx.URLParams.Add("id", testKeyID)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/keys/"+testKeyID, nil)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		handler.ServeHTTP(w, r)
		err := json.Unmarshal(w.Body.Bytes(), &keyResponse)

		assert.Nil(t, err)
		assert.IsType(t, key{}, keyResponse)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetKeyHandler_ServeHTTP_Negative(t *testing.T) {
	// creating storage mocks
	keyGetter := &mock.Getter{KeyError: storage.ErrKeyNotFound}
	// creating handler
	handler := NewGetKeyHandler(keyGetter)
	// Populate the request's context with our test data.
	rctx := chi.NewRouteContext()

	t.Run("Getting a key without providing an ID", func(t *testing.T) {
		var errResponse api.ErrResponse

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/key", nil)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		handler.ServeHTTP(w, r)
		err := json.Unmarshal(w.Body.Bytes(), &errResponse)

		assert.Nil(t, err)
		assert.IsType(t, api.ErrResponse{}, errResponse)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestGetKeyHandler_ServeHTTP_Failure(t *testing.T) {
	// creating storage mocks
	keyGetter := &mock.Getter{KeyError: errors.New("failed fetch data from DB")}
	// creating handler
	handler := NewGetKeyHandler(keyGetter)
	// Populate the request's context with our test data.
	rctx := chi.NewRouteContext()

	t.Run("Getting a key ", func(t *testing.T) {
		var errResponse api.ErrResponse
		testKeyID := "test-bla-bla"

		rctx.URLParams.Add("id", testKeyID)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/keys/"+testKeyID, nil)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		handler.ServeHTTP(w, r)
		err := json.Unmarshal(w.Body.Bytes(), &errResponse)

		assert.Nil(t, err)
		assert.IsType(t, api.ErrResponse{}, errResponse)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
