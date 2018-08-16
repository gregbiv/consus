package model

import (
	"time"

	"github.com/bluele/factory-go/factory"
	"github.com/lib/pq"
)

// KeyFactory represents factory method of Key
var KeyFactory = factory.NewFactory(
	&Key{
		KeyID:     "test-key-id",
		Value:     "some random value",
		CreatedAt: time.Now(),
		UpdatedAt: pq.NullTime{Time: time.Now(), Valid: true},
		ExpiresAt: pq.NullTime{Valid: false},
	},
)

// Key is the mapping to the keys database table.
type Key struct {
	KeyID     string      `db:"id"`
	Value     string      `db:"value"`
	CreatedAt time.Time   `db:"created_at"`
	UpdatedAt pq.NullTime `db:"updated_at"`
	ExpiresAt pq.NullTime `db:"expires_at"`
}
