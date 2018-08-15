package keys

import (
	"time"

	"github.com/gregbiv/sandbox/pkg/model"
)

// key describes an API model
type key struct {
	KeyID     string     `json:"id"`
	Value     string     `json:"value"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	ExpiresAt *time.Time `json:"expires_at"`
}

func (s *key) fromDB(dbKey *model.Key) {
	s.KeyID = dbKey.KeyID
	s.CreatedAt = dbKey.CreatedAt
	s.Value = dbKey.Value

	if dbKey.UpdatedAt.Valid {
		updatedAt := dbKey.UpdatedAt.Time
		s.UpdatedAt = &updatedAt
	}

	if dbKey.ExpiresAt.Valid {
		expiresAt := dbKey.ExpiresAt.Time
		s.ExpiresAt = &expiresAt
	}
}
