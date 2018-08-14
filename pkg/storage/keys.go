package storage

import "github.com/gregbiv/sandbox/pkg/model"

type (
	// Getter is responsible for SELECTing keys
	Getter interface {
		// GetByID gets an addOn subscription model from DB by ID
		GetAll() ([]*model.Key, error)
	}
)
