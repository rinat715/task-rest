package sqlstore

import (
	"database/sql"
	"go_rest/internal/logger"
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
	logger.Info("DB close")
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
