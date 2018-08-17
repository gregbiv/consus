package command

import (
	"strings"

	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database"
	"github.com/mattes/migrate/source/file"
	log "github.com/sirupsen/logrus"

	// Enables Postgres for migrations
	_ "github.com/mattes/migrate/database/postgres"
	// Enables stub DB for testing
	_ "github.com/mattes/migrate/database/stub"
)

// MigrateCommand represents a CLI migrate command
type MigrateCommand struct {
	Meta
}

// Run execute the migration command
func (c *MigrateCommand) Run(args []string) int {
	flags := c.FlagSet("http")
	flags.Usage = func() { c.UI.Output(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	if c.Config.Migration.Dir == "" {
		log.Errorf("empty option: dir. args: %v", args)
		return 1
	}

	if c.Config.Migration.Version == 0 {
		log.Errorf("empty option: -version. args: %v", args)
		return 1
	}

	log.Infof("trying to migrate DB to version %d", c.Config.Migration.Version)

	db, err := database.Open(c.Config.Database.PostgresDB.DSN)
	if err != nil {
		log.Fatalf("Cannot initialize db with the following DSN: %s", c.Config.Database.PostgresDB.DSN)
		return 1
	}

	f := &file.File{}
	sourceDriver, err := f.Open("file://" + c.Config.Migration.Dir)
	if err != nil {
		log.Error(err)
		return 1
	}

	// Configure migration
	migration, err := migrate.NewWithInstance("file", sourceDriver, "database", db)

	if err != nil {
		log.Error(err)
		return 1
	}

	// Migrate the system to the correct version
	err = migration.Migrate(c.Config.Migration.Version)
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Infof("Database already at version %d", c.Config.Migration.Version)
			return 0
		}

		log.Errorf("Database migration error: %s", err)
		return 1
	}

	log.Infof("Migrated database to version %d", c.Config.Migration.Version)
	return 0
}

// Help outputs a helper text for the command
func (*MigrateCommand) Help() string {
	helpText := `
Usage: consus migrate [options]

  Migrate database to the given version
`

	return strings.TrimSpace(helpText)
}

// Synopsis of the migrate command
func (c *MigrateCommand) Synopsis() string {
	return "Migrate database version"
}
