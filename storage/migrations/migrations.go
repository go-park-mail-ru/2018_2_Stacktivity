package migrations

import (
	"2018_2_Stacktivity/storage"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

func InitMigration() {
	driver, err := postgres.WithInstance(storage.GetInstance(), &postgres.Config{})
	if err != nil {
		return
	}

	migration, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return
	}

	migration.Up()
}
