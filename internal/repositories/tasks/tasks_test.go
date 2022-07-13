package tasks

import (
	"database/sql"
	"testing"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/mattn/go-sqlite3"
)

func MakeMigrate(t *testing.T, db *sql.DB) {
	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		t.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"root/db/migration/",
		"sqlite3",
		instance)
	if err != nil {
		t.Fatal(err)
	}
	m.Up()
}

// TODO все переделать
func TestDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	MakeMigrate(t, db)
	ts := NewTaskRepository(db)

	stmt, err := db.Prepare("INSERT INTO task(text, date, done) VALUES(?, ?, ?)")
	if err != nil {
		t.Fatal(err)
	}

	date1, err := time.Parse("2006-01-02", "2009-10-23")
	res, err := stmt.Exec("Example text", date1, false)
	if err != nil {
		t.Fatal(err)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	stmt.Close()

	stmt, err = db.Prepare("INSERT INTO tag(value, taskid) VALUES(?, ?)")
	if err != nil {
		t.Fatal(err)
	}
	res, err = stmt.Exec("Text tag1", lastId)
	if err != nil {
		t.Fatal(err)
	}

	res, err = stmt.Exec("Text tag2", lastId)
	if err != nil {
		t.Fatal(err)
	}
	stmt.Close()

	task, err := ts.Get(int(lastId))
	if err != nil {
		t.Fatal(err)
	}
	if task.Text != "Example text" {
		t.Errorf("Error: got %v, want %v", task.Text, "Example text")
	}
	if len(task.Tags) != 3 {
		t.Errorf("Error: task.Tags != 2")
	}
}
