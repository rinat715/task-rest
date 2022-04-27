package sqlstore

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
	*tasks
}

func New(db *sql.DB) (*Store, error) {
	return &Store{
		db,
		&tasks{
			db,
			&tag{db},
		},
	}, nil
}

func (ts *tasks) Close() error {
	return ts.db.Close()
}

func (ts *tasks) CreateTables() error {
	_, err := ts.db.Exec(
		`CREATE TABLE task (
			taskid INTEGER NOT NULL,
			text VARCHAR(255) NOT NULL,
			date Date NOT NULL,
			done BOOLEAN DEFAULT false,
			PRIMARY KEY (taskid)
		);
		
		CREATE TABLE tag (
			tagid INT,
			value VARCHAR(255),
			taskid INT,
			PRIMARY KEY (tagid),
			FOREIGN KEY (taskid) REFERENCES task(taskid)
		);`,
	)
	return err
}
