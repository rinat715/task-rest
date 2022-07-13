package sqlitedb

import (
	"context"
	"database/sql"

	"go_rest/internal/config"

	_ "github.com/mattn/go-sqlite3"
)

func NewConnDB(ctx context.Context, c *config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", c.PathDB)
	if err != nil {
		return &sql.DB{}, err
	}
	go func() {
		<-ctx.Done()
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()
	return db, nil
}
