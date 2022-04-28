// Eli Bendersky [https://eli.thegreenplace.net]
// This code is in the public domain.
package store

import (
	m "go_rest/internal/models"
	"testing"
	"time"
)

func getIdcreatedTask(ts *TaskStore) int {
	taskCreated := m.Task{
		Text: "Hola",
		Tags: nil,
		Date: m.JsonDate(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)),
		Done: false,
	}
	ts.Create(&taskCreated)
	return taskCreated.Id
}

func TestCreateAndGet(t *testing.T) {
	// Create a store and a single task.
	ts := New()
	id := getIdcreatedTask(ts)

	// We should be able to retrieve this task by ID, but nothing with other
	// IDs.
	task, err := ts.Get(id)
	if err != nil {
		t.Fatal(err)
	}

	if task.Id != id {
		t.Errorf("got task.Id=%d, id=%d", task.Id, id)
	}
	if task.Text != "Hola" {
		t.Errorf("got Text=%v, want %v", task.Text, "Hola")
	}

	// Asking for all tasks, we only get the one we put in.
	allTasks := ts.GetAll()
	if len(allTasks) != 1 || allTasks[0].Id != id {
		t.Errorf("got len(allTasks)=%d, allTasks[0].Id=%d; want 1, %d", len(allTasks), allTasks[0].Id, id)
	}

	_, err = ts.Get(id + 1)
	if err == nil {
		t.Fatal("got nil, want error")
	}

	// Add another task. Expect to find two tasks in the store.
	ts.Create(m.Task{
		Text: "hey",
		Tags: nil,
		Date: m.JsonDate(time.Now()),
		Done: false,
	})
	allTasks2 := ts.GetAll()
	if len(allTasks2) != 2 {
		t.Errorf("got len(allTasks2)=%d; want 2", len(allTasks2))
	}
}

func TestDelete(t *testing.T) {
	ts := New()
	id1 := ts.Create(m.Task{
		Text: "Foo",
		Tags: nil,
		Date: JsonDate(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)),
		Done: false,
	}).Id
	id2 := ts.Create(Task{
		Text: "Bar",
		Tags: nil,
		Date: JsonDate(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)),
		Done: false,
	}).Id

	if err := ts.Delete(id1 + 1001); err == nil {
		t.Fatalf("delete task id=%d, got no error; want error", id1+1001)
	}

	if err := ts.Delete(id1); err != nil {
		t.Fatal(err)
	}
	if err := ts.Delete(id1); err == nil {
		t.Fatalf("delete task id=%d, got no error; want error", id1)
	}

	if err := ts.Delete(id2); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteAll(t *testing.T) {
	ts := New()
	ts.Create(Task{
		Text: "Foo",
		Tags: nil,
		Date: JsonDate(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)),
		Done: false,
	})
	ts.Create(Task{
		Text: "Bar",
		Tags: nil,
		Date: JsonDate(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)),
		Done: false,
	})

	if err := ts.DeleteAll(); err != nil {
		t.Fatal(err)
	}

	tasks := ts.GetAll()
	if len(tasks) > 0 {
		t.Fatalf("want no tasks remaining; got %v", tasks)
	}
}

func TestGetTasksByTag(t *testing.T) {
	ts := New()
	ts.Create(Task{
		Text: "XY",
		Tags: []string{"Movies"},
		Date: JsonDate(time.Now()),
		Done: false,
	})

	ts.Create(Task{
		Text: "XY",
		Tags: []string{"Bills"},
		Date: JsonDate(time.Now()),
		Done: false,
	})

	ts.Create(Task{
		Text: "XY",
		Tags: []string{"Bills"},
		Date: JsonDate(time.Now()),
		Done: false,
	})

	ts.Create(Task{
		Text: "YZR",
		Tags: []string{"Bills", "Movies"},
		Date: JsonDate(time.Now()),
		Done: false,
	})

	ts.Create(Task{
		Text: "YWZ",
		Tags: []string{"Movies", "Bills"},
		Date: JsonDate(time.Now()),
		Done: false,
	})

	ts.Create(Task{
		Text: "WZT",
		Tags: []string{"Movies"},
		Date: JsonDate(time.Now()),
		Done: false,
	})

	var tests = []struct {
		tag     string
		wantNum int
	}{
		{"Movies", 3},
		{"Bills", 4},
		{"Ferrets", 0},
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			numByTag := len(ts.GetByTag(tt.tag))

			if numByTag != tt.wantNum {
				t.Errorf("got %v, want %v", numByTag, tt.wantNum)
			}
		})
	}
}

func TestGetTasksByDueDate(t *testing.T) {
	mustParseDate := func(tstr string) JsonDate {
		tt, err := JsonDateParse(tstr)
		if err != nil {
			t.Fatal(err)
		}
		return tt
	}

	ts := New()
	ts.Create(Task{
		Text: "XY1",
		Tags: nil,
		Date: mustParseDate("2020-12-01"),
		Done: false,
	})

	ts.Create(Task{
		Text: "XY1",
		Tags: nil,
		Date: mustParseDate("2000-12-21"),
		Done: false,
	})

	ts.Create(Task{
		Text: "XY1",
		Tags: nil,
		Date: mustParseDate("2020-12-01"),
		Done: false,
	})

	ts.Create(Task{
		Text: "XY1",
		Tags: nil,
		Date: mustParseDate("2000-12-21"),
		Done: false,
	})

	ts.Create(Task{
		Text: "XY1",
		Tags: nil,
		Date: mustParseDate("1991-01-01"),
		Done: false,
	})

	// Check a single task can be fetched.
	d, _ := JsonDateParse("1991-01-01")
	tasks1 := ts.GetByDate(d)
	if len(tasks1) != 1 {
		t.Errorf("got len=%d, want 1", len(tasks1))
	}
	if tasks1[0].Text != "XY5" {
		t.Errorf("got Text=%s, want XY5", tasks1[0].Text)
	}

	var tests = []struct {
		date    string
		wantNum int
	}{
		{"2020-01-01", 0},
		{"2020-12-01", 2},
		{"2000-12-21", 2},
		{"1991-01-01", 1},
		{"2020-12-21", 0},
	}

	for _, tt := range tests {
		t.Run(tt.date, func(t *testing.T) {
			date := mustParseDate(tt.date)
			numByDate := len(ts.GetByDate(date))

			if numByDate != tt.wantNum {
				t.Errorf("got %v, want %v", numByDate, tt.wantNum)
			}
		})
	}
}
