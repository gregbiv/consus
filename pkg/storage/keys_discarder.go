package storage

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/palantir/stacktrace"
)

type dbDiscarder struct {
	db *sqlx.DB
}

// NewDiscarder inits and returns a Discarder instance
func NewDiscarder(db *sqlx.DB) Discarder {
	return &dbDiscarder{db: db}
}

func (dd *dbDiscarder) Discard(ID string) error {
	tx, err := dd.db.Beginx()
	if err != nil {
		return err
	}

	if err := discardKey(tx, ID); err != nil {
		if err == ErrKeyNotFound {
			return err
		}
		log.Panic(stacktrace.Propagate(err, "Failed to discard key", err, ID))
	}

	if err := tx.Commit(); err != nil {
		log.Panic(stacktrace.Propagate(err, "Failed to commit transaction to delete a key (Key ID: %s)", err, ID))
	}

	return nil
}

func (dd *dbDiscarder) Truncate() error {
	tx, err := dd.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`TRUNCATE keys CASCADE `)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Panic(stacktrace.Propagate(err, "Failed to commit truncate transaction", err))
	}

	return nil
}

func discardKey(tx *sqlx.Tx, ID string) error {
	query := `
        DELETE FROM keys
		WHERE id = $1
    `

	result, err := tx.Exec(query, ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrKeyNotFound
	}

	return nil
}
