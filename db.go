package main

import (
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
	log "github.com/sirupsen/logrus"
)

// DB is the shared database connection pool
var DB *pop.Connection

func initDB() error {
	// pop requires a database.yml file
	// This yaml file refers to the DATABASE_URL environment variable as the uri
	// So we set this env variable, so that database.yml can point to it.
	// I hope this changes in the future
	if envy.Get("DATABASE_URI", "") != "" {
		envy.Set("DATABASE_URI", settings.Database.URI)
	}

	var err error
	DB, err = pop.Connect(settings.Environment)

	if err != nil {
		return err
	}

	log.Info("Established DB connection.")
	pop.Debug = settings.Environment == development
	return nil
}

func migrateDB() error {
	fileMigrator, err := pop.NewFileMigrator("./migrations", DB)

	if err != nil {
		return err
	}

	fileMigrator.Status()

	return fileMigrator.Up()
}
