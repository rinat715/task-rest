package sqlitestore

import (
	"fmt"
	m "go_rest/internal/models"
	s "go_rest/internal/sqlstore"
	"reflect"
	"testing"
	"time"
)

var task m.Task = m.Task{
	Text: "Example text",
	Tags: []string{
		"test_tag1",
		"test_tag2",
	},
	Date: m.JsonDate(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)),
	Done: false,
}

func createDb() (*s.Store, error) {
	ts, err := New(":memory:")
	if err != nil {
		return ts, err
	}
	ts.CreateTables()
	if err != nil {
		return ts, err
	}
	return ts, nil
}

func createTask(t *testing.T, ts *s.Store) m.Task {
	err := ts.Create(&task)
	if err != nil {
		t.Fatal(err)
	}
	return task
}

func TestCreate(t *testing.T) {
	var err error
	ts, err := createDb()
	if err != nil {
		t.Fatal(err)
	}
	task := createTask(t, ts)

	if task.Id != 1 {
		t.Errorf("Error: got %v, want %v", task.Id, 1)
	}

	defer ts.Close()
}

func TestGet(t *testing.T) {
	var err error
	ts, err := createDb()
	if err != nil {
		t.Fatal(err)
	}
	createTask(t, ts)

	taskfromDb, err := ts.Get(1)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(taskfromDb, task) {
		t.Errorf("Error: got %v", taskfromDb)
	}

	_, err = ts.Get(2)
	if err == nil {
		t.Error("Not exception raised")
	}
	if fmt.Sprint(err) != "sql: no rows in result set" {
		t.Errorf("Unknown error message %v", err)
	}
	defer ts.Close()
}
