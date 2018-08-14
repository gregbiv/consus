package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Specification represents structured configuration variables
type Specification struct {
	Debug    bool   `envconfig:"DEBUG" default:"false"`
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
	Port     int    `envconfig:"PORT" default:"8090"`
	Database struct {
		PostgresDB struct {
			DSN string `envconfig:"DATABASE_DSN"`
		}
	}
}

// LoadEnv loads config variables into Specification
func LoadEnv() *Specification {
	var conf Specification
	envconfig.MustProcess("", &conf)

	return &conf
}
