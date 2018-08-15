package storage

import (
	"fmt"
	"strings"
	"time"

	"github.com/gregbiv/sandbox/pkg/model"
	"github.com/jmoiron/sqlx"
	"github.com/palantir/stacktrace"
)

type dbKeyStorer struct {
	db *sqlx.DB
}

// NewStorer inits and returns a KeyStorer instance
func NewStorer(db *sqlx.DB) Storer {
	return &dbKeyStorer{db: db}
}

func (m *dbKeyStorer) Store(d *model.Key) error {
	tx, err := m.db.Beginx()
	if err != nil {
		return err
	}

	if err := m.insertKey(tx, d); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (m *dbKeyStorer) insertKey(tx *sqlx.Tx, key *model.Key) error {
	params := map[string]interface{}{
		"id":         key.KeyID,
		"value":      key.Value,
		"created_at": time.Now(),
	}

	keys := make([]string, 0, len(params))
	named := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
		named = append(named, fmt.Sprintf(":%s", key))
	}

	query := fmt.Sprintf(`INSERT INTO keys (%s) VALUES (%s)`, strings.Join(keys, ", "), strings.Join(named, ", "))

	nstmt, err := tx.PrepareNamed(query)

	if err != nil {
		return stacktrace.Propagate(
			err,
			"failed to create a prepared statement to store data in key table",
			err.Error())
	}

	defer nstmt.Close()

	if _, err = nstmt.Exec(params); err != nil {
		return stacktrace.Propagate(err, "failed to store data into key table", err.Error())
	}

	return nil
}
