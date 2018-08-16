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

func TestNewUpdater(t *testing.T) {
	db, _, _ := dbMock.NewDbMock()
	defer db.Close()

	expectedResult := &dbKeyUpdater{
		db: db,
	}

	result := NewUpdater(db)

	assert.Equal(t, expectedResult, result)
}

func TestKeysUpdater(t *testing.T) {
	t.Parallel()

	ID := "test"
	value := "random value"

	y := 1990
	m := time.Month(12)
	d := 1
	createdAt := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	updatedAt := pq.NullTime{Valid: false}
	expiresAt := pq.NullTime{Time: time.Now(), Valid: true}

	key := &model.Key{
		KeyID:     ID,
		Value:     value,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		ExpiresAt: expiresAt,
	}

	query := `^UPDATE keys SET  \(.*\) =  \(.*\) WHERE id = \$1`

	t.Run("Updates a key", func(t *testing.T) {
		// Arrange
		db, dbmock, _ := dbMock.NewDbMock()
		defer db.Close()
		st := NewUpdater(db)
		execResult := sqlmock.NewResult(int64(1), int64(1))

		dbmock.ExpectBegin()
		dbmock.ExpectPrepare(query).ExpectExec().WithArgs().
			WillReturnResult(execResult)

		dbmock.ExpectCommit()

		// Act
		err := st.Update(key)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Returns an error when unable to begin a transaction", func(t *testing.T) {
		db, dbmock, _ := dbMock.NewDbMock()
		defer db.Close()
		st := NewUpdater(db)
		expectedErr := errors.New("unable to begin a transaction")
		dbmock.ExpectBegin().WillReturnError(expectedErr)

		err := st.Update(key)

		assert.EqualError(t, err, expectedErr.Error())
	})

	t.Run("Rolls back and returns an error when exec fails", func(t *testing.T) {
		// Arrange
		db, dbmock, _ := dbMock.NewDbMock()
		defer db.Close()
		st := NewUpdater(db)
		expectedErr := errors.New("transaction failed")

		dbmock.ExpectBegin()
		dbmock.ExpectPrepare(query).ExpectExec().WillReturnError(expectedErr)
		dbmock.ExpectRollback()

		// Act
		err := st.Update(key)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectedErr.Error())
	})

	t.Run("Rolls back and returns an error when prepare fails", func(t *testing.T) {
		// Arrange
		db, dbmock, _ := dbMock.NewDbMock()
		defer db.Close()
		st := NewUpdater(db)
		expectedErr := errors.New("transaction failed")

		dbmock.ExpectBegin()
		dbmock.ExpectPrepare(query).WillReturnError(expectedErr)
		dbmock.ExpectRollback()

		// Act
		err := st.Update(key)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectedErr.Error())
	})

	t.Run("Rolls back and returns an error when no rows were affected", func(t *testing.T) {
		// Arrange
		db, dbmock, _ := dbMock.NewDbMock()
		defer db.Close()
		st := NewUpdater(db)
		execResult := sqlmock.NewResult(int64(0), int64(0))

		dbmock.ExpectBegin()

		dbmock.ExpectPrepare(query).ExpectExec().WithArgs().
			WillReturnResult(execResult)

		dbmock.ExpectRollback()

		// Act
		err := st.Update(key)

		// Assert
		assert.EqualError(t, err, ErrUpdatingKey.Error())
	})

	t.Run("Rolls back and returns an error when unable to fetch rows affected", func(t *testing.T) {
		// Arrange
		db, dbmock, _ := dbMock.NewDbMock()
		defer db.Close()
		st := NewUpdater(db)
		expectedErr := errors.New("transaction failed")
		execResult := sqlmock.NewErrorResult(expectedErr)
		dbmock.ExpectBegin()
		dbmock.ExpectPrepare(query).ExpectExec().WithArgs().
			WillReturnResult(execResult)
		dbmock.ExpectRollback()

		// Act
		err := st.Update(key)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectedErr.Error())
	})
}
