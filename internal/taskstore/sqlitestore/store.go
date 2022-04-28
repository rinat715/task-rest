package sqlitestore

import (
	"database/sql"

	s "go_rest/internal/taskstore/sqlstore"

	_ "github.com/mattn/go-sqlite3"
)

func New(path string) (*s.Store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return &s.Store{}, err
	}
	return s.New(db)
}
