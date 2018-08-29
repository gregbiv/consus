package keys

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/gregbiv/consus/pkg/api"
	"github.com/gregbiv/consus/pkg/storage"
	"github.com/gregbiv/consus/pkg/storage/mock"
	"github.com/stretchr/testify/assert"
)

func TestDiscardKeyHandler_ServeHTTP(t *testing.T) {
	t.Parallel()

	// creating storage mocks
	keyDiscarder := &mock.Discarder{}
	// creating handler
	handler := NewDiscardKeyHandler(keyDiscarder)

	// Populate the request's context with our test data.
	rctx := chi.NewRouteContext()

	t.Run("Discarding a key", func(t *testing.T) {
		testKeyID := "test-bla-bla"
		rctx.URLParams.Add("id", testKeyID)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("DELETE", "/keys/"+testKeyID, nil)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		handler.ServeHTTP(w, r)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}

func TestDiscardKeyHandler_ServeHTTP_Negative(t *testing.T) {
	// creating storage mocks
	keyDiscarder := &mock.Discarder{DiscardError: storage.ErrKeyNotFound}
	// creating handler
	handler := NewDiscardKeyHandler(keyDiscarder)

	// Populate the request's context with our test data.
	rctx := chi.NewRouteContext()

	t.Run("Discarding a key", func(t *testing.T) {
		var errResponse api.ErrResponse
		testKeyID := "test-bla-bla"
		rctx.URLParams.Add("id", testKeyID)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("DELETE", "/keys/"+testKeyID, nil)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		handler.ServeHTTP(w, r)
		err := json.Unmarshal(w.Body.Bytes(), &errResponse)

		assert.Nil(t, err)
		assert.IsType(t, api.ErrResponse{}, errResponse)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestDiscardKeyHandler_ServeHTTP_FailureDiscarder(t *testing.T) {
	// creating storage mocks
	keyDiscarder := &mock.Discarder{DiscardError: errors.New("failed remove data from DB")}
	// creating handler
	handler := NewDiscardKeyHandler(keyDiscarder)

	// Populate the request's context with our test data.
	rctx := chi.NewRouteContext()

	t.Run("Discarding a key with error in DB", func(t *testing.T) {
		testKeyID := "test-bla-bla"
		rctx.URLParams.Add("id", testKeyID)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("DELETE", "/keys/"+testKeyID, nil)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		handler.ServeHTTP(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
