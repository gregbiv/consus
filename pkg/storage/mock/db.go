package mock

import (
	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

// NewDbMock creates a mock for sqlx.DB and DB interaction mock
func NewDbMock() (*sqlx.DB, sqlmock.Sqlmock, error) {
	mockDB, mock, err := sqlmock.New()

	return sqlx.NewDb(mockDB, "sqlmock"), mock, err
}
