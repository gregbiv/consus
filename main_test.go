package main

import (
	"github.com/DATA-DOG/godog"
	"github.com/gregbiv/sandbox/features/bootstrap"
	"github.com/gregbiv/sandbox/pkg/config"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

var (
	URL = "http://localhost:8090"
)

func FeatureContext(s *godog.Suite) {
	cfg := config.LoadEnv()

	bootstrap.RegisterGomega(s)

	conn, err := sqlx.Open("postgres", cfg.Database.PostgresDB.DSN)
	if err != nil {
		log.Fatal(err)
	}

	bootstrap.RegisterSystemContext(s, URL)

	bootstrap.RegisterKeysContext(s, URL, conn)
}
