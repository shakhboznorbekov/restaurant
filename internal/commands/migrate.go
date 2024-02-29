package commands

import (
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/restaurant/internal/pkg/repository/postgresql"
	"log"
)

// ErrHelp provides context that help was given.
var ErrHelp = errors.New("provided help")

func MigrateUp(db *postgresql.Database) {
	driver, err := postgres.WithInstance(db.DB.DB, &postgres.Config{})
	if err != nil {
		log.Fatal("error in opening driver: ", err.Error())
	}
	m, err := migrate.NewWithDatabaseInstance("file://internal/pkg/scripts", "restaurant", driver)
	if err != nil {
		log.Fatal("error in creating migrations: ", err.Error())
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			er := dbErrorMsg(db, err.Error())
			if er != nil {
				log.Fatalf("Error in writing actual error to database [ %v ]", er)
			}
		}
	}
}

func dbErrorMsg(db *postgresql.Database, err string) error {
	query := fmt.Sprintf(`UPDATE schema_migrations SET error='/*[%s]*/'`, err)

	_, er := db.Exec(query)
	if er != nil {
		return er
	}

	return nil
}
