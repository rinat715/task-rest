package sqlstore

import (
	"database/sql"
	m "go_rest/internal/models"
	"reflect"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestTag(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	ts, err := New(db)
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

	var tags = make(map[int][]string)

	tags[1] = []string{
		"test_tag1",
		"test_tag2",
	}

	tags[2] = []string{
		"test_tag3",
		"test_tag4",
	}

	err = ts.Create(&m.Task{
		Text: "Example text",
		Tags: tags[1],
		Date: m.JsonDate(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)),
		Done: false,
	})
	if err != nil {
		t.Error(err)
	}

	err = ts.Create(&m.Task{
		Text: "Example text2",
		Tags: tags[2],
		Date: m.JsonDate(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)),
		Done: false,
	})
	if err != nil {
		t.Error(err)
	}

	t.Run("Get tags by taskId", func(t *testing.T) {
		res, err := ts.getTagsbyTaskid(1)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(res, tags[1]) {
			t.Errorf("unknown tags: %s", res)
		}
		if reflect.DeepEqual(res, tags[2]) {
			t.Errorf("unknown tags: %s", res)
		}
	})

	t.Run("Get tags by taskIds", func(t *testing.T) {
		res, err := ts.getTagsbyTaskids([]int{1, 2})
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(res[1], tags[1]) {
			t.Errorf("unknown tags: %v", res)
		}
		if !reflect.DeepEqual(res[2], tags[2]) {
			t.Errorf("unknown tags: %v", res)
		}
	})

	t.Run("Get all tags", func(t *testing.T) {
		res, err := ts.getTags()
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(res[1], tags[1]) {
			t.Errorf("unknown tags: %v", res)
		}
		if !reflect.DeepEqual(res[2], tags[2]) {
			t.Errorf("unknown tags: %v", res)
		}
	})

	t.Cleanup(func() {
		ts.Close()
	})
}
