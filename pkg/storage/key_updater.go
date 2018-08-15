package storage

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gregbiv/sandbox/pkg/model"
	"github.com/jmoiron/sqlx"
	"github.com/palantir/stacktrace"
)

// ErrUpdatingKey failed to update a key
var ErrUpdatingKey = errors.New("failed to update a key")

type dbKeyUpdater struct {
	db *sqlx.DB
}

// NewUpdater inits and returns a KeyUpdater instance
func NewUpdater(db *sqlx.DB) Updater {
	return &dbKeyUpdater{db: db}
}

func (du *dbKeyUpdater) Update(key *model.Key) error {
	tx, err := du.db.Beginx()
	if err != nil {
		return err
	}

	steps := []func(*sqlx.Tx, *model.Key) (bool, error){
		du.updateKey,
	}

	for _, step := range steps {
		ok, err := step(tx, key)
		if err != nil {
			tx.Rollback()
			stacktrace.Propagate(err, ErrUpdatingKey.Error(), err.Error())
			return err
		}

		if !ok {
			tx.Rollback()
			return ErrUpdatingKey
		}
	}

	return tx.Commit()
}

func (du *dbKeyUpdater) updateKey(tx *sqlx.Tx, key *model.Key) (bool, error) {
	params := map[string]interface{}{
		"id":         key.KeyID,
		"value":      key.Value,
		"updated_at": time.Now(),
		"expires_at": nil,
	}

	if key.ExpiresAt.Valid {
		params["expires_at"] = key.ExpiresAt.Time
	}

	keys := make([]string, 0, len(params))
	named := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
		named = append(named, fmt.Sprintf(":%s", key))
	}

	query := fmt.Sprintf(`UPDATE keys SET (%s) = (%s) WHERE id = :id`, strings.Join(keys, ", "), strings.Join(named, ", "))
	nstmt, err := tx.PrepareNamed(query)

	if err != nil {
		return false, stacktrace.Propagate(
			err,
			"failed to create a prepared statement to store data in key table",
			err.Error(), query, params)
	}

	result, err := du.executeQuery(nstmt, params)
	nstmt.Close()

	return result, err
}

func (du *dbKeyUpdater) executeQuery(ns *sqlx.NamedStmt, params interface{}) (bool, error) {
	result, err := ns.Exec(params)
	if err != nil {
		return false, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
