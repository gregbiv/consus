package storage

import (
	"errors"
	"testing"

	dbMock "github.com/gregbiv/consus/pkg/storage/mock"
	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestNewDiscarder(t *testing.T) {
	db, _, _ := dbMock.NewDbMock()
	defer db.Close()
	expectedResult := &dbDiscarder{
		db: db,
	}

	result := NewDiscarder(db)

	assert.Equal(t, expectedResult, result)
}

func TestKeyManager_Discard(t *testing.T) {
	t.Parallel()

	ID := "test"
	query := `^DELETE FROM keys WHERE id = \$1`

	t.Run("Fails when transaction fails to begin", func(t *testing.T) {
		// Arrange
		sqlxdb, dbmock, _ := dbMock.NewDbMock()
		defer sqlxdb.Close()
		discarder := NewDiscarder(sqlxdb)
		expectedErr := errors.New("failed to begin a transaction")
		dbmock.ExpectBegin().WillReturnError(expectedErr)

		// Act
		err := discarder.Discard(ID)

		// Assert
		assert.EqualError(t, err, expectedErr.Error())
	})

	t.Run("Discards a key", func(t *testing.T) {
		// Arrange
		sqlxdb, dbmock, _ := dbMock.NewDbMock()
		defer sqlxdb.Close()
		res := sqlmock.NewResult(1, 1)
		discarder := NewDiscarder(sqlxdb)
		dbmock.ExpectBegin()
		dbmock.ExpectExec(query).WithArgs(
			sqlmock.AnyArg(),
		).WillReturnResult(res)
		dbmock.ExpectCommit()

		// Act
		err := discarder.Discard(ID)

		// Assert
		assert.Nil(t, err)
	})

	t.Run("Panics when transaction fails to commit", func(t *testing.T) {
		// Arrange
		sqlxdb, dbmock, _ := dbMock.NewDbMock()
		defer sqlxdb.Close()
		res := sqlmock.NewResult(1, 1)
		discarder := NewDiscarder(sqlxdb)
		dbmock.ExpectBegin()
		dbmock.ExpectExec(query).WithArgs(
			sqlmock.AnyArg(),
		).WillReturnResult(res)
		expectedErr := errors.New("failed to commit a transaction")
		dbmock.ExpectCommit().WillReturnError(expectedErr)

		// Act
		result := func() {
			discarder.Discard(ID)
		}

		// Assert
		assert.Panics(t, result)
	})

	t.Run("Panics when unable to discard", func(t *testing.T) {
		// Arrange
		sqlxdb, dbmock, _ := dbMock.NewDbMock()
		defer sqlxdb.Close()
		discarder := NewDiscarder(sqlxdb)
		dbmock.ExpectBegin()
		expectedErr := errors.New("failed discard a key")
		dbmock.ExpectExec(query).WithArgs(
			sqlmock.AnyArg(),
		).WillReturnError(expectedErr)
		dbmock.ExpectCommit().WillReturnError(expectedErr)

		// Act
		result := func() {
			discarder.Discard(ID)
		}

		// Assert
		assert.Panics(t, result)
	})

	t.Run("Returns an error when unable to find a key to discard",
		func(t *testing.T) {
			// Arrange
			sqlxdb, dbmock, _ := dbMock.NewDbMock()
			defer sqlxdb.Close()
			discarder := NewDiscarder(sqlxdb)
			dbmock.ExpectBegin()
			expectedErr := ErrKeyNotFound
			dbmock.ExpectExec(query).WithArgs(
				sqlmock.AnyArg(),
			).WillReturnError(expectedErr)
			dbmock.ExpectCommit().WillReturnError(expectedErr)

			// Act
			err := discarder.Discard(ID)

			// Assert
			assert.EqualError(t, err, expectedErr.Error())
		})

	t.Run("Returns an error when no rows affected", func(t *testing.T) {
		// Arrange
		sqlxdb, dbmock, _ := dbMock.NewDbMock()
		defer sqlxdb.Close()
		discarder := NewDiscarder(sqlxdb)
		dbmock.ExpectBegin()
		res := sqlmock.NewResult(0, 0)
		dbmock.ExpectExec(query).WithArgs(
			sqlmock.AnyArg(),
		).WillReturnResult(res)
		dbmock.ExpectCommit()

		// Act
		err := discarder.Discard(ID)

		// Assert
		assert.EqualError(t, err, ErrKeyNotFound.Error())
	})

	t.Run("Returns an error when unable to fetch rows affected", func(t *testing.T) {
		// Arrange
		sqlxdb, dbmock, _ := dbMock.NewDbMock()
		defer sqlxdb.Close()
		discarder := NewDiscarder(sqlxdb)
		dbmock.ExpectBegin()
		expectedErr := errors.New("transaction failed")
		res := sqlmock.NewErrorResult(expectedErr)
		dbmock.ExpectExec(query).WithArgs(
			sqlmock.AnyArg(),
		).WillReturnResult(res)
		dbmock.ExpectCommit()

		// Act
		result := func() { discarder.Discard(ID) }

		// Assert
		assert.Panics(t, result)
	})
}

func TestKeyManager_Truncate(t *testing.T) {
	t.Parallel()

	query := `^TRUNCATE keys CASCADE`

	t.Run("Fails when transaction fails to begin", func(t *testing.T) {
		// Arrange
		sqlxdb, dbmock, _ := dbMock.NewDbMock()
		defer sqlxdb.Close()
		discarder := NewDiscarder(sqlxdb)
		expectedErr := errors.New("failed to begin a transaction")
		dbmock.ExpectBegin().WillReturnError(expectedErr)

		// Act
		err := discarder.Truncate()

		// Assert
		assert.EqualError(t, err, expectedErr.Error())
	})

	t.Run("Discards all keys", func(t *testing.T) {
		// Arrange
		sqlxdb, dbmock, _ := dbMock.NewDbMock()
		defer sqlxdb.Close()
		res := sqlmock.NewResult(1, 1)
		discarder := NewDiscarder(sqlxdb)
		dbmock.ExpectBegin()
		dbmock.ExpectExec(query).WillReturnResult(res)
		dbmock.ExpectCommit()

		// Act
		err := discarder.Truncate()

		// Assert
		assert.Nil(t, err)
	})

	t.Run("Panics when transaction fails to commit", func(t *testing.T) {
		// Arrange
		sqlxdb, dbmock, _ := dbMock.NewDbMock()
		defer sqlxdb.Close()
		res := sqlmock.NewResult(1, 1)
		discarder := NewDiscarder(sqlxdb)
		dbmock.ExpectBegin()
		dbmock.ExpectExec(query).WillReturnResult(res)
		expectedErr := errors.New("failed to commit a transaction")
		dbmock.ExpectCommit().WillReturnError(expectedErr)

		// Act
		result := func() {
			discarder.Truncate()
		}

		// Assert
		assert.Panics(t, result)
	})

	t.Run("Panics when unable to truncate", func(t *testing.T) {
		// Arrange
		sqlxdb, dbmock, _ := dbMock.NewDbMock()
		defer sqlxdb.Close()
		discarder := NewDiscarder(sqlxdb)
		dbmock.ExpectBegin()
		expectedErr := errors.New("failed to truncate")
		dbmock.ExpectExec(query).WillReturnError(expectedErr)
		dbmock.ExpectCommit().WillReturnError(expectedErr)

		// Act
		result := func() {
			discarder.Truncate()
		}

		// Assert
		assert.Panics(t, result)
	})
}
