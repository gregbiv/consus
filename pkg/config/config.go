package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Specification represents structured configuration variables
type Specification struct {
	Debug     bool   `envconfig:"DEBUG" default:"false"`
	LogLevel  string `envconfig:"LOG_LEVEL" default:"info"`
	Port      int    `envconfig:"PORT" default:"8090"`
	Migration Migration
	Database  struct {
		PostgresDB struct {
			DSN string `envconfig:"DATABASE_DSN"`
		}
	}
}

// Migration config
type Migration struct {
	Version uint   `envconfig:"DATABASE_VERSION" default:"1"`
	Dir     string `envconfig:"MIGRATION_DIR" default:"resources/migrations"`
}

// LoadEnv loads config variables into Specification
func LoadEnv() *Specification {
	var conf Specification
	envconfig.MustProcess("", &conf)

	return &conf
}
