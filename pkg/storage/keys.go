package storage

import (
	"errors"

	"github.com/gregbiv/consus/pkg/model"
)

var (
	// ErrKeyNotFound ...
	ErrKeyNotFound = errors.New("key not found")
)

type (
	// Getter is responsible for SELECTing keys
	Getter interface {
		// GetAll gets all keys
		GetAll(filterStr string, activeOnly bool) ([]*model.Key, error)
		// GetByID gets a key by ID
		GetByID(ID string, activeOnly bool) (*model.Key, error)
	}

	// Discarder is responsible for DELETing keys
	Discarder interface {
		// Discard removes entry from database
		Discard(ID string) error
		// Truncate removes everything from database
		Truncate() error
	}

	// Storer is responsible for storing keys
	Storer interface {
		// Store stores a key in database
		Store(key *model.Key) error
	}

	// Updater is responsible for updating keys
	Updater interface {
		// Update updates a key in database
		Update(key *model.Key) error
	}
)
