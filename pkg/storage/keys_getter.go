package storage

import (
	"database/sql"

	"github.com/gregbiv/sandbox/pkg/model"
	"github.com/jmoiron/sqlx"
)

type dbGetter struct {
	db *sqlx.DB
}

// NewGetter inits and returns a Getter instance
func NewGetter(db *sqlx.DB) Getter {
	return &dbGetter{db: db}
}

func (dg *dbGetter) GetAll() ([]*model.Key, error) {
	var list []*model.Key
	var err error

	query := `
		SELECT
			id, 
			value,
			created_at, 
			updated_at, 
			expires_at
		FROM keys
	`

	err = dg.db.Select(&list, query)

	if err != nil {
		if err == sql.ErrNoRows {
			return list, nil
		}

		return nil, err
	}

	return list, nil
}
