package keys

import (
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gregbiv/consus/pkg/model"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestToModel(t *testing.T) {
	id := "test"
	value := "some random value"
	timestamp := time.Now()

	key := &key{
		KeyID:     id,
		Value:     value,
		CreatedAt: timestamp,
		UpdatedAt: &timestamp,
		ExpiresAt: &timestamp,
	}

	result, err := key.toModel()

	assert.Nil(t, err)
	assert.Equal(t, id, result.KeyID)
	assert.Equal(t, value, result.Value)
	assert.IsType(t, time.Time{}, result.CreatedAt)
	assert.IsType(t, pq.NullTime{}, result.UpdatedAt)
	assert.IsType(t, pq.NullTime{}, result.ExpiresAt)
}

func TestToModel_WithInvalidID(t *testing.T) {
	id := "this id is not valid"
	value := "some random value"
	timestamp := time.Now()

	key := &key{
		KeyID:     id,
		Value:     value,
		CreatedAt: timestamp,
		UpdatedAt: &timestamp,
		ExpiresAt: &timestamp,
	}

	_, err := key.toModel()

	assert.Error(t, err)
	assert.IsType(t, ErrInvalidKeyID, err)
}

func TestFromDb(t *testing.T) {
	result := key{}
	id := "test"
	value := "some random value"
	timestamp := time.Now()
	modelKey := &model.Key{
		KeyID:     id,
		Value:     value,
		CreatedAt: timestamp,
		UpdatedAt: pq.NullTime{Time: timestamp, Valid: true},
		ExpiresAt: pq.NullTime{Time: timestamp, Valid: true},
	}

	result.fromDB(modelKey)

	assert.Equal(t, id, result.KeyID)
	assert.Equal(t, value, result.Value)
}

func TestFromRequest(t *testing.T) {
	result := key{}
	payload := `
	{
		"id" : "test",
		"value" : "some random value"
	}
	`
	req := httptest.NewRequest("PUT", "/keys", strings.NewReader(payload))

	err := result.fromRequest(req)

	assert.Nil(t, err)
	assert.Equal(t, "test", result.KeyID)
	assert.Equal(t, "some random value", result.Value)
}

func TestFromRequestWithInvalidBody(t *testing.T) {
	result := key{}
	req := httptest.NewRequest("PUT", "/keys", strings.NewReader("{ broken body"))

	err := result.fromRequest(req)

	assert.IsType(t, ErrInvalidBody, err)
}

func TestFromRequestWithBrokenReader(t *testing.T) {
	result := key{}
	file, err := os.Open("my_file.zip")
	req := httptest.NewRequest("PUT", "/keys", file)
	err = result.fromRequest(req)

	assert.Error(t, err)
}
