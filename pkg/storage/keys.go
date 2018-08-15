package storage

import (
	"errors"

	"github.com/gregbiv/sandbox/pkg/model"
)

var (
	// ErrKeyNotFound ...
	ErrKeyNotFound = errors.New("key not found")
)

type (
	// Getter is responsible for SELECTing keys
	Getter interface {
		// GetAll gets all keys
		GetAll() ([]*model.Key, error)
		// GetByID gets a key by ID
		GetByID(ID string) (*model.Key, error)
	}

	// Discarder is responsible for DELETing keys
	Discarder interface {
		// Discard removes entry from database
		Discard(ID string) error
		// Truncate removes everything from database
		Truncate() error
	}
)
