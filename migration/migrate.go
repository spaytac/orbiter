package migration

import (
	"embed"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/sirupsen/logrus"
	"os"
)

//go:embed schema/*.sql
var fs embed.FS

func Migrate() {
	d, err := iofs.New(fs, "schema")
	if err != nil {
		panic(fmt.Sprintf("unable to initiate migration due to: %s", err.Error()))
	}
	DbPath := os.Getenv("DB_PATH")
	if DbPath == "" {
		panic("unable to initiate migration due to: DB_PATH env not set")
	}

	m, err := migrate.NewWithSourceInstance(
		"iofs",
		d,
		fmt.Sprintf("sqlite://%s", DbPath),
	)
	if err != nil {
		panic(fmt.Sprintf("unable to initiate migration due to: %s", err.Error()))
	}

	previousVersion, _, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			logrus.Info("Previous migration version not found.")
			previousVersion = 0
		} else {
			panic(fmt.Sprintf("unable to get current version of schema migration: %s", err.Error()))
		}
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(fmt.Sprintf("unable to run migration due to: %s", err.Error()))
	}

	if err == nil {
		currentVersion, _, err := m.Version()
		if err != nil {
			panic(fmt.Sprintf("unable to get current version of schema migration: %s", err.Error()))
		}

		logrus.WithField("prev:", previousVersion).WithField("curr:", currentVersion).Info("Migration complete")
	} else {
		logrus.Info("No migration pending")
	}

	if sourceErr, dbErr := m.Close(); sourceErr != nil || dbErr != nil {
		logrus.WithError(sourceErr).WithError(dbErr).Info("Error when closing source and database connections.")
	}
}
