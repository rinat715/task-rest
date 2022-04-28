package sqlitestore

import (
	"errors"
	e "go_rest/internal/errors"
	m "go_rest/internal/models"
	s "go_rest/internal/taskstore/sqlstore"
	"reflect"
	"testing"
)

func TestDb(t *testing.T) {
	var err error
	date1, err := m.JsonDateParse("2009-10-23")
	if err != nil {
		t.Fatal(err)
	}

	var task m.Task = m.Task{
		Text: "Example text",
		Tags: []string{
			"test_tag1",
			"test_tag2",
		},
		Date: date1,
		Done: false,
	}

	ts, err := New(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	ts.CreateTables()
	if err != nil {
		t.Fatal(err)
	}

	if err != nil {
		t.Fatal(err)
	}

	err = ts.Create(&task)
	if err != nil {
		t.Error(err)
	}

	date2, err := m.JsonDateParse("2010-10-23")
	if err != nil {
		t.Fatal(err)
	}

	err = ts.Create(&m.Task{
		Text: "Example text2",
		Tags: []string{
			"test_tag2",
		},
		Date: date2,
		Done: false,
	})
	if err != nil {
		t.Error(err)
	}

	t.Run("Create task", func(t *testing.T) {
		if task.Id != 1 {
			t.Errorf("Error: got %v, want %v", task.Id, 1)
		}
	})
	t.Run("Get task", func(t *testing.T) {
		taskfromDb, err := ts.Get(1)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(taskfromDb, task) {
			t.Errorf("Error: got %v", taskfromDb)
		}
		if task.Text != "Example text" {
			t.Errorf("Error: got %v, want %v", task.Text, "Example text")
		}

		_, err = ts.Get(3)
		if err == nil {
			t.Errorf("Error TaskNotFound not raised")
		} else {
			var ew e.ErrorWrapper
			if errors.As(err, &ew) {
				if !errors.Is(ew.Err, &s.TaskNotFound{Id: 3}) {
					t.Error(ew.Err)
				}
			} else {
				t.Error(err)
			}
		}

	})
	t.Run("Get all task", func(t *testing.T) {
		res, err := ts.GetAll()
		if err != nil {
			t.Error(err)
		}
		if len(res) != 2 {
			t.Errorf("Error: length tasks not equal 2")
		}
	})
	t.Run("Get task by tag", func(t *testing.T) {
		res, err := ts.GetByTag("test_tag1")
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(res[0], task) {
			t.Errorf("Error: got %v", res[0])
		}

		res, err = ts.GetByTag("test_tag2")
		if err != nil {
			t.Error(err)
		}
		if len(res) != 2 {
			t.Errorf("Error: length tasks not equal 2")
		}
		_, err = ts.GetByTag("test_tag3")
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Get task by date", func(t *testing.T) {
		date, err := m.JsonDateParse("2009-10-23")
		if err != nil {
			t.Fatal(err)
		}

		res, err := ts.GetByDate(date)
		if err != nil {
			t.Error(err)
		}
		if len(res) == 0 {
			t.Fatal("Not found tasks")
		}
		if !reflect.DeepEqual(res[0], task) {
			t.Errorf("Error: got %v", res[0])
		}
	})

	t.Cleanup(func() {
		ts.Close()
	})
}
