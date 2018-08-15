package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/satori/go.uuid"
)

type (
	// URLExtractor interface for extracting uuid/country/date from the request
	URLExtractor interface {
		// ExtractUUIDFromRoute extracts the uuid from the url
		UUIDFromRoute(r *http.Request, routeParam string) (*uuid.UUID, error)
	}

	urlExtractor struct{}
)

// NewURLExtractor is the constructor for the UUID Extractor interface
func NewURLExtractor() URLExtractor {
	return &urlExtractor{}
}

func (e *urlExtractor) UUIDFromRoute(r *http.Request, routeParam string) (*uuid.UUID, error) {
	routeUUID, err := uuid.FromString(chi.URLParam(r, routeParam))
	if err != nil {
		return nil, err
	}

	return &routeUUID, nil
}
