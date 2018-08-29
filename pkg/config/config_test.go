package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	setGlobalConfigEnv()

	cfg := LoadEnv()
	assertConfig(t, cfg)
}

func assertConfig(t *testing.T, cfg *Specification) {
	assert.Equal(t, false, cfg.Debug)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, 8090, cfg.Port)
	assert.Equal(t, "postgres://", cfg.Database.PostgresDB.DSN)
	assert.Equal(t, uint(1), cfg.Migration.Version)
	assert.Equal(t, "resources/migrations", cfg.Migration.Dir)
}

func setGlobalConfigEnv() {
	os.Setenv("DEBUG", "false")
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("PORT", "8090")
	os.Setenv("DATABASE_DSN", "postgres://")
	os.Setenv("DATABASE_VERSION", "1")
	os.Setenv("MIGRATION_DIR", "resources/migrations")
}
