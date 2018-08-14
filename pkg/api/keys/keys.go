package keys

import (
	"time"

	"github.com/gregbiv/sandbox/pkg/model"
	"github.com/satori/go.uuid"
)

// key describes an API model
type key struct {
	KeyID     *uuid.UUID `json:"key_id"`
	Value     string     `json:"value"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	ExpiresAt *time.Time `json:"expires_at"`
}

func (s *key) fromDB(dbKey *model.Key) error {
	id, err := uuid.FromString(dbKey.KeyID)
	if err != nil {
		return err
	}

	s.KeyID = &id
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

	return nil
}
