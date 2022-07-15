package repositories

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	_ "github.com/golang-migrate/migrate/source/file"
)

const migrationPath string = "file://./../migration"

func makeMigrate(db *sql.DB) error {
	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"sqlite",
		instance)
	if err != nil {
		return err
	}
	return m.Up()
}

func SetupTestdb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return &sql.DB{}, err
	}
	err = makeMigrate(db)
	if err != nil {
		return &sql.DB{}, err
	}
	return db, nil
}
