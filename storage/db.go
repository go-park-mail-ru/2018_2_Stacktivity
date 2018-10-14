package storage

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var db *sqlx.DB

func InitDB(DSN string) error {
	var err error
	// TODO add Docker for db
	db, err = sqlx.Connect("postgres", DSN)
	if err != nil {
		return errors.Wrap(err, "can't open database")
	}
	err = db.Ping()
	if err != nil {
		return errors.Wrap(err, "can't connect to database")
	}
	return nil
}

func GetInstance() *sql.DB {
	return db.DB
}
