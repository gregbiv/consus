package keys

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gregbiv/sandbox/pkg/storage/mock"
	"github.com/stretchr/testify/assert"
)

func TestDiscardKeysHandler_ServeHTTP(t *testing.T) {
	t.Parallel()

	// creating storage mocks
	keyDiscarder := &mock.Discarder{}
	// creating handler
	handler := NewDiscardKeysHandler(keyDiscarder)

	t.Run("Discarding all keys", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("DELETE", "/keys", nil)

		handler.ServeHTTP(w, r)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}

func TestDiscardKeysHandler_ServeHTTP_FailureDiscarder(t *testing.T) {
	// creating storage mocks
	keyDiscarder := &mock.Discarder{DiscardError: errors.New("failed remove data from DB")}
	// creating handler
	handler := NewDiscardKeysHandler(keyDiscarder)

	t.Run("Discarding a subscription key with error in DB", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("DELETE", "/keys", nil)

		handler.ServeHTTP(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
