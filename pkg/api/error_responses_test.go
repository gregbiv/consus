package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorResponses_NotFound(t *testing.T) {
	r := httptest.NewRequest("GET", "/smth", nil)
	w := httptest.NewRecorder()

	NotFound(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `
		{
			"error": {
				"code": "InvalidUri",
				"message": "The requested URI does not represent any resource on the server."
			}
		}
	`, w.Body.String())
}

func TestErrorResponses_RenderErrMissingURIParam(t *testing.T) {
	r := httptest.NewRequest("GET", "/smth", nil)
	w := httptest.NewRecorder()
	missingParam := "interval_id"

	RenderErrMissingURIParam(w, r, missingParam)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `
		{
			"error": {
				"code": "MissingUriParam",
				"message": "The 'interval_id' query parameter is required."
			}
		}
	`, w.Body.String())
}

func TestErrorResponses_RenderInternalServerError(t *testing.T) {
	r := httptest.NewRequest("GET", "/smth", nil)
	w := httptest.NewRecorder()
	err := fmt.Errorf("Internal server error")

	RenderInternalServerError(w, r, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `
		{
			"error": {
				"code": "InternalError",
				"message": "Internal server error"
			}
		}
	`, w.Body.String())
}

func TestErrorResponses_RenderInvalidInput(t *testing.T) {
	r := httptest.NewRequest("GET", "/smth", nil)
	w := httptest.NewRecorder()
	invalidParam := "interval_type"
	errorMessage := "field is invalid"

	RenderInvalidInput(w, r, invalidParam, errorMessage)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `
		{
			"error": {
				"code": "InvalidInput",
				"message": "field is invalid",
				"target": "interval_type"
			}
		}
	`, w.Body.String())
}

func TestErrorResponses_RenderBadGateway(t *testing.T) {
	r := httptest.NewRequest("GET", "/smth", nil)
	w := httptest.NewRecorder()
	err := fmt.Errorf("Bad gateway")

	RenderBadGateway(w, r, err)

	assert.Equal(t, http.StatusBadGateway, w.Code)
	assert.JSONEq(t, `
		{
			"error": {
				"code": "BadGateway",
				"message": "Bad gateway"
			}
		}
	`, w.Body.String())
}
