package model

import (
	"time"

	"github.com/lib/pq"
)

// Key is the mapping to the keys database table.
type Key struct {
	KeyID     string      `db:"id"`
	Value     string      `db:"value"`
	CreatedAt time.Time   `db:"created_at"`
	UpdatedAt pq.NullTime `db:"updated_at"`
	ExpiresAt pq.NullTime `db:"expires_at"`
}
