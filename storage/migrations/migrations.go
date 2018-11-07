package migrations

import (
	"2018_2_Stacktivity/storage"

	"log"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

func InitMigration() {
	driver, err := postgres.WithInstance(storage.GetInstance(), &postgres.Config{})
	if err != nil {
		log.Println("can't get instanse ", err)
		return
	}

	migration, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		log.Println("can't migrage: ", err)
		return
	}

	if err := migration.Up(); err != nil {
		log.Println("can't up migrations: ", err)
		return
	}
}
