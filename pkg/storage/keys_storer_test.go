package storage

import (
	"errors"
	"testing"
	"time"

	"github.com/gregbiv/sandbox/pkg/model"
	dbMock "github.com/gregbiv/sandbox/pkg/storage/mock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestNewStorer(t *testing.T) {
	db, _, _ := dbMock.NewDbMock()
	defer db.Close()

	expectedResult := &dbKeyStorer{
		db: db,
	}

	result := NewStorer(db)

	assert.Equal(t, expectedResult, result)
}

func TestKeysStorer(t *testing.T) {
	t.Parallel()

	ID := "test"
	value := "random value"

	y := 1990
	m := time.Month(12)
	d := 1
	createdAt := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	updatedAt := pq.NullTime{Valid: false}
	expiresAt := pq.NullTime{Valid: false}

	key := &model.Key{
		KeyID:     ID,
		Value:     value,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		ExpiresAt: expiresAt,
	}

	query := `^INSERT INTO keys \(.*\) VALUES \(.*\)`

	t.Run("Creates a key", func(t *testing.T) {
		// Arrange
		db, dbmock, _ := dbMock.NewDbMock()
		defer db.Close()
		st := NewStorer(db)
		execResult := sqlmock.NewResult(int64(1), int64(1))

		dbmock.ExpectBegin()
		dbmock.ExpectPrepare(query).ExpectExec().WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg()).WillReturnResult(execResult)

		dbmock.ExpectCommit()

		// Act
		err := st.Store(key)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Returns an error when unable to begin a transaction", func(t *testing.T) {
		db, dbmock, _ := dbMock.NewDbMock()
		defer db.Close()
		st := NewStorer(db)
		expectedErr := errors.New("Unable to begin a transaction")
		dbmock.ExpectBegin().WillReturnError(expectedErr)

		err := st.Store(key)

		assert.EqualError(t, err, expectedErr.Error())
	})

	t.Run("Rolls back and returns an error when exec fails", func(t *testing.T) {
		// Arrange
		db, dbmock, _ := dbMock.NewDbMock()
		defer db.Close()
		st := NewStorer(db)
		expectedErr := errors.New("Transaction failed")

		dbmock.ExpectBegin()
		dbmock.ExpectPrepare(query).ExpectExec().WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg()).WillReturnError(expectedErr)

		dbmock.ExpectRollback()

		// Act
		err := st.Store(key)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectedErr.Error())
	})

	t.Run("Rolls back and returns an error when prepare fails", func(t *testing.T) {
		// Arrange
		db, dbmock, _ := dbMock.NewDbMock()
		defer db.Close()
		st := NewStorer(db)
		expectedErr := errors.New("Transaction failed")

		dbmock.ExpectBegin()
		dbmock.ExpectPrepare(query).WillReturnError(expectedErr)

		dbmock.ExpectRollback()

		// Act
		err := st.Store(key)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectedErr.Error())
	})
}
