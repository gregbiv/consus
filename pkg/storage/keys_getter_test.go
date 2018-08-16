package storage

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/gregbiv/sandbox/pkg/model"
	dbMock "github.com/gregbiv/sandbox/pkg/storage/mock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestNewGetter(t *testing.T) {
	db, _, _ := dbMock.NewDbMock()
	expectedResult := &dbGetter{
		db: db,
	}

	result := NewGetter(db)

	assert.Equal(t, expectedResult, result)
}

func TestDbGetter_GetByID(t *testing.T) {
	t.Parallel()

	sqlxdb, dbmock, _ := dbMock.NewDbMock()
	defer sqlxdb.Close()
	dbGetter := NewGetter(sqlxdb)

	expectedColumns := []string{
		"id",
		"value",
		"created_at",
		"updated_at",
		"expires_at",
	}

	ID := "test"
	value := "random value"

	query := `SELECT \* FROM keys WHERE id = \$1 AND \(NOT expires_at < NOW\(\) OR expires_at IS NULL\)`

	y := 1990
	m := time.Month(12)
	d := 1
	createdAt := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	updatedAt := pq.NullTime{Valid: false}
	expiresAt := pq.NullTime{Valid: false}

	t.Run("Gets a key by ID", func(t *testing.T) {
		// Arrange
		expectedResult := &model.Key{
			KeyID:     ID,
			Value:     value,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			ExpiresAt: expiresAt,
		}

		expectedRows := sqlmock.NewRows(expectedColumns).AddRow(
			ID,
			value,
			createdAt,
			updatedAt,
			expiresAt,
		)

		dbmock.ExpectQuery(query).WillReturnRows(expectedRows)

		// Act
		result, err := dbGetter.GetByID(ID, true)

		// Assert
		assert.NoError(t, err)
		assert.EqualValues(t, expectedResult, result)
	})

	t.Run("Returns key not found when no rows in DB", func(t *testing.T) {
		// Arrange
		expectedRows := sqlmock.NewRows(expectedColumns)
		dbmock.ExpectQuery(query).WillReturnRows(expectedRows)

		// Act
		result, err := dbGetter.GetByID(ID, true)

		// Assert
		assert.EqualError(t, err, ErrKeyNotFound.Error())
		assert.Nil(t, result)
	})

	t.Run("Returns an error when table does not exists", func(t *testing.T) {
		// Arrange
		expectedError := errors.New("table not found")
		dbmock.ExpectQuery(query).WillReturnError(expectedError)

		// Act
		result, err := dbGetter.GetByID(ID, true)

		// Assert
		assert.EqualError(t, err, expectedError.Error())
		assert.Nil(t, result)
	})
}

func TestDbGetter_GetAll(t *testing.T) {
	t.Parallel()

	sqlxdb, dbmock, _ := dbMock.NewDbMock()
	defer sqlxdb.Close()
	dbGetter := NewGetter(sqlxdb)

	expectedColumns := []string{
		"id",
		"value",
		"created_at",
		"updated_at",
		"expires_at",
	}

	ID := "test"
	value := "random value"
	filterStr := "test"

	query := `SELECT \* FROM keys WHERE \(NOT expires_at < NOW\(\) OR expires_at IS NULL\)` +
		` AND LIKE '%' || $1 || '%'`

	y := 1990
	m := time.Month(12)
	d := 1
	createdAt := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	updatedAt := pq.NullTime{Valid: false}
	expiresAt := pq.NullTime{Valid: false}

	t.Run("Gets all keys", func(t *testing.T) {
		// Arrange
		expectedResult := []*model.Key{
			{
				KeyID:     ID,
				Value:     value,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				ExpiresAt: expiresAt,
			},
		}

		expectedRows := sqlmock.NewRows(expectedColumns).AddRow(
			ID,
			value,
			createdAt,
			updatedAt,
			expiresAt,
		)

		dbmock.ExpectQuery(query).WillReturnRows(expectedRows)

		// Act
		result, err := dbGetter.GetAll(filterStr, true)

		// Assert
		assert.NoError(t, err)
		assert.EqualValues(t, expectedResult, result)
	})

	t.Run("Returns key not found when no rows in DB", func(t *testing.T) {
		// Arrange
		var expectedResult []*model.Key
		expectedError := sql.ErrNoRows
		dbmock.ExpectQuery(query).WillReturnError(expectedError)

		// Act
		result, err := dbGetter.GetAll(filterStr, true)

		// Assert
		assert.NoError(t, err)
		assert.EqualValues(t, expectedResult, result)
	})

	t.Run("Returns an error when table does not exists", func(t *testing.T) {
		// Arrange
		expectedError := errors.New("table not found")
		dbmock.ExpectQuery(query).WillReturnError(expectedError)

		// Act
		result, err := dbGetter.GetAll(filterStr, true)

		// Assert
		assert.EqualError(t, err, expectedError.Error())
		assert.Nil(t, result)
	})
}
