package main

import (
	"os"

	"github.com/gregbiv/consus/pkg/command"
	"github.com/gregbiv/consus/pkg/config"
	"github.com/mitchellh/cli"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
)

func init() {
	prometheus.MustRegister(version.NewCollector("service"))
}

func main() {
	log.SetFormatter(&log.JSONFormatter{FieldMap: log.FieldMap{
		log.FieldKeyTime: "@timestamp",
		log.FieldKeyMsg:  "message",
	}})

	// Config
	conf := config.LoadEnv()

	// Logging
	level, err := log.ParseLevel(conf.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	log.SetLevel(level)
	if conf.Debug {
		log.Debugf("Initialized with config: %+v", conf)
	}

	c := &cli.CLI{
		Name:     "consus",
		Version:  "dev",
		HelpFunc: cli.BasicHelpFunc("consus"),
		Commands: commands(conf),
		Args:     os.Args[1:],
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}

func commands(conf *config.Specification) map[string]cli.CommandFactory {
	meta := command.Meta{
		UI: &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
		Config: conf,
	}

	cf := map[string]cli.CommandFactory{
		"http": func() (cli.Command, error) {
			return &command.HTTPCommand{
				Meta: meta,
			}, nil
		},
		"migrate": func() (cli.Command, error) {
			return &command.MigrateCommand{
				Meta: meta,
			}, nil
		},
	}

	return cf
}
