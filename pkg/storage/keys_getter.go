package storage

import (
	"database/sql"
	"fmt"
	"strings"

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

func (dg *dbGetter) GetAll(filterStr string, activeOnly bool) ([]*model.Key, error) {
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

	if activeOnly {
		query += ` WHERE (NOT expires_at < NOW() OR expires_at IS NULL)`
	}

	if filterStr != "" {
		filterStr = strings.Replace(filterStr, "$", "%", -1)
		query += fmt.Sprintf(" AND id LIKE '%%' || '%s' || '%%'", filterStr)
	}

	err = dg.db.Select(&list, query)

	if err != nil {
		if err == sql.ErrNoRows {
			return list, nil
		}

		return nil, err
	}

	return list, nil
}

func (dg *dbGetter) GetByID(ID string, activeOnly bool) (*model.Key, error) {
	var key model.Key
	var err error

	query := `
		SELECT
			id, 
			value,
			created_at, 
			updated_at, 
			expires_at
		FROM keys
		WHERE id = $1
	`

	if activeOnly {
		query += ` AND NOT expires_at < NOW() OR expires_at IS NULL`
	}

	err = dg.db.Get(&key, query, ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrKeyNotFound
		}

		return nil, err
	}

	return &key, nil
}
