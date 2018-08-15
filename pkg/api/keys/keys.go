package keys

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/gregbiv/sandbox/pkg/model"
	"github.com/lib/pq"
)

// basic regular expression for validating strings
const Alpha = "^[a-zA-Z]+$"

var (
	// ErrInvalidBody represents the error when the request body is invalid.
	ErrInvalidBody = errors.New("invalid request body provided")
	// ErrInvalidKey represents the error when the provided key is invalid
	ErrInvalidKeyID = errors.New("invalid id provided, value should contain only letters")
)

// key describes an API model
type key struct {
	KeyID     string     `json:"id"`
	Value     string     `json:"value"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	ExpiresAt *time.Time `json:"expires_at"`
}

func (k *key) fromDB(dbKey *model.Key) {
	k.KeyID = dbKey.KeyID
	k.CreatedAt = dbKey.CreatedAt
	k.Value = dbKey.Value

	if dbKey.UpdatedAt.Valid {
		updatedAt := dbKey.UpdatedAt.Time
		k.UpdatedAt = &updatedAt
	}

	if dbKey.ExpiresAt.Valid {
		expiresAt := dbKey.ExpiresAt.Time
		k.ExpiresAt = &expiresAt
	}
}

func (k *key) toModel() (modelKey model.Key, err error) {
	// validating key ID
	if !regexp.MustCompile(Alpha).MatchString(k.KeyID) {
		return modelKey, ErrInvalidKeyID
	}

	if k.ExpiresAt != nil {
		expiresAt := k.ExpiresAt
		modelKey.ExpiresAt = pq.NullTime{Time: *expiresAt, Valid: !expiresAt.IsZero()}
	}

	modelKey.KeyID = k.KeyID
	modelKey.Value = k.Value

	return modelKey, nil
}

func (k *key) fromRequest(r *http.Request) error {
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, k)
	if err != nil {
		return ErrInvalidBody
	}

	return nil
}
