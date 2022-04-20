// Eli Bendersky [https://eli.thegreenplace.net]
// This code is in the public domain.
package taskstore

import (
	"testing"
	"time"
)

func TestCreateAndGet(t *testing.T) {
	// Create a store and a single task.
	ts := New()
	var d JsonDate = JsonDate(time.Now())
	id := ts.Create("Hola", nil, d, false).Id

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
	ts.Create("hey", nil, JsonDate(time.Now()), false)
	allTasks2 := ts.GetAll()
	if len(allTasks2) != 2 {
		t.Errorf("got len(allTasks2)=%d; want 2", len(allTasks2))
	}
}

func TestDelete(t *testing.T) {
	ts := New()
	id1 := ts.Create("Foo", nil, JsonDate(time.Now()), false).Id
	id2 := ts.Create("Bar", nil, JsonDate(time.Now()), false).Id

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
	ts.Create("Foo", nil, JsonDate(time.Now()), false)
	ts.Create("Bar", nil, JsonDate(time.Now()), false)

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
	ts.Create("XY", []string{"Movies"}, JsonDate(time.Now()), false)
	ts.Create("YZ", []string{"Bills"}, JsonDate(time.Now()), false)
	ts.Create("YZR", []string{"Bills"}, JsonDate(time.Now()), false)
	ts.Create("YWZ", []string{"Bills", "Movies"}, JsonDate(time.Now()), false)
	ts.Create("WZT", []string{"Movies", "Bills"}, JsonDate(time.Now()), false)

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
	timeFormat := "2006-01-02"
	mustParseDate := func(tstr string) JsonDate {
		tt, err := time.Parse(timeFormat, tstr)
		if err != nil {
			t.Fatal(err)
		}
		return JsonDate(tt)
	}

	ts := New()
	ts.Create("XY1", nil, mustParseDate("2020-12-01"), false)
	ts.Create("XY2", nil, mustParseDate("2000-12-21"), false)
	ts.Create("XY3", nil, mustParseDate("2020-12-01"), false)
	ts.Create("XY4", nil, mustParseDate("2000-12-21"), false)
	ts.Create("XY5", nil, mustParseDate("1991-01-01"), false)

	// Check a single task can be fetched.
	y, m, d := time.Time(mustParseDate("1991-01-01")).Date()
	tasks1 := ts.GetByDate(y, m, d)
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
			y, m, d := time.Time(mustParseDate(tt.date)).Date()
			numByDate := len(ts.GetByDate(y, m, d))

			if numByDate != tt.wantNum {
				t.Errorf("got %v, want %v", numByDate, tt.wantNum)
			}
		})
	}
}
