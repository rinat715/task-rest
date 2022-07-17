package tasks

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	m "go_rest/internal/models"
	r "go_rest/internal/repositories"
)

func createUser(db *sql.DB) (int, error) {
	stmt, err := db.Prepare("INSERT INTO users(email, pass, is_admin) VALUES(?, ?, ?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec("test@email.com", "123", true)
	if err != nil {
		return 0, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	return int(lastId), nil
}

func setup(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO tasks(text, date, done, userid) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}

	date1, _ := time.Parse("2006-01-02", "2009-10-23")
	_, err = stmt.Exec("Example text", date1, false, 1)
	if err != nil {
		return err
	}

	date2, _ := time.Parse("2006-01-02", "2009-10-22")
	_, err = stmt.Exec("Example text2", date2, true, 1)
	if err != nil {
		return err
	}

	stmt.Close()

	stmt, err = db.Prepare("INSERT INTO tags(value, taskid) VALUES(?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec("Text tag1", 1)
	if err != nil {
		return err
	}

	_, err = stmt.Exec("Text tag2", 1)
	if err != nil {
		return err
	}

	_, err = stmt.Exec("Text tag3", 2)
	if err != nil {
		return err
	}
	stmt.Close()
	return nil
}

func TestDb(t *testing.T) {
	db, err := r.SetupTestdb()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	_, err = createUser(db)
	if err != nil {
		t.Fatal(err)
	}

	setup(db)

	ts := NewTaskRepository(db)

	t.Run("Get task", func(t *testing.T) {

		task, err := ts.Get(1)
		if err != nil {
			t.Fatal(err)
		}
		if task.Text != "Example text" {
			t.Errorf("Error: got %v, want %v", task.Text, "Example text")
		}
		if len(task.Tags) != 2 {
			t.Errorf("Error: task.Tags != 2")
		}
		task, _ = ts.Get(3)
		if !task.IsEmpty() {
			t.Errorf("non empty task")
		}
	})

	t.Run("Get all task", func(t *testing.T) {
		tasks, err := ts.GetAll()
		if err != nil {
			t.Fatal(err)
		}
		if len(tasks) != 2 {
			t.Log(len(tasks))
			t.Errorf("Error: tasks != 2")
		}
	})

	t.Run("Get by tag", func(t *testing.T) {
		tasks, err := ts.GetByTag("Text tag2")
		if err != nil {
			t.Fatal(err)
		}
		if len(tasks) != 1 {
			t.Errorf("Error: tasks != 1")
		} else {
			task := tasks[0]
			if task.Text != "Example text" {
				t.Errorf("Error неправильный текст таски")
			}
		}

	})
	t.Run("Get by date", func(t *testing.T) {
		// TODO не работает
		tasks, err := ts.GetByDate("2009-10-23T00:00:00Z")
		if err != nil {
			t.Fatal(err)
		}
		if len(tasks) != 1 {
			t.Log(tasks)
			t.Errorf("Error: tasks != 1")
		} else {
			task := tasks[0]
			if task.Text != "Example text2" {
				t.Errorf("Error неправильный текст таски")
			}
		}
	})
	t.Run("Delete by id", func(t *testing.T) {
		err := ts.Delete(1)
		if err != nil {
			t.Fatal(err)
		}
		var deletedId int
		row := db.QueryRow("SELECT taskid FROM tasks WHERE taskid = ?", 1)
		err = row.Scan(&deletedId)
		if err != sql.ErrNoRows {
			t.Log(deletedId)
			t.Errorf("ошибка удаления")
		}
		var countTags int
		db.QueryRow("SELECT COUNT(tagid) FROM tags WHERE taskid = ?", 1).Scan(&countTags)
		if countTags != 0 {
			t.Log(countTags)
			t.Errorf("Ошибка при удалении тегов")
		}
	})
}

func TestDeleteAll(t *testing.T) {
	db, err := r.SetupTestdb()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	_, err = createUser(db)
	if err != nil {
		t.Fatal(err)
	}

	setup(db)

	ts := NewTaskRepository(db)
	err = ts.DeleteAll()
	if err != nil {
		t.Fatal(err)
	}
	var countTags, countTasks int
	db.QueryRow("SELECT COUNT(tagid) FROM tags").Scan(&countTags)
	db.QueryRow("SELECT COUNT(taskid) FROM tasks").Scan(&countTags)
	if countTags != 0 || countTasks != 0 {
		t.Log(countTags)
		t.Log(countTasks)
		t.Errorf("Ошибка при удалении")
	}
}

func TestCreate(t *testing.T) {
	db, err := r.SetupTestdb()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	_, err = createUser(db)
	if err != nil {
		t.Fatal(err)
	}

	setup(db)

	ts := NewTaskRepository(db)

	var tags = make(m.Tags, 2)

	tags[0] = m.Tag{
		Text: "Created tag1",
	}

	tags[1] = m.Tag{
		Text: "Created tag2",
	}
	date, _ := time.Parse("2006-01-02", "2022-07-01")

	task := m.Task{
		UserId: 1,
		Text:   "Created task",
		Tags:   tags,
		Date:   date,
		Done:   false,
	}

	err = ts.Create(&task, 1)
	if err != nil {
		t.Fatal(err)
	}

	if task.Id != 3 {
		t.Errorf("Error: got %v, want %v", task.Id, 3)
	}
	var taskId int
	var text, dateGet string
	var done bool
	db.QueryRow("SELECT taskid, text, date, done FROM tasks WHERE taskid = ?", 3).Scan(&taskId, &text, &dateGet, &done)
	if taskId != 3 {
		t.Errorf("Error: got %v, want %v", taskId, 3)
	}
	if text != task.Text {
		t.Errorf("Error: got %v, want %v", text, task.Text)
	}
	if dateGet != "2022-07-01T00:00:00Z" {
		t.Errorf("Error: got %v, want %v", dateGet, "2022-07-01T00:00:00Z")
	}
	if done != false {
		t.Errorf("Error: got %v, want %v", done, false)
	}

	var countTags int
	db.QueryRow("SELECT COUNT(*) FROM tags WHERE taskid = ?", 3).Scan(&countTags)
	if countTags != 2 {
		t.Log(countTags)
		t.Errorf("Ошибка при создании тегов")
	}

	var tagId int
	var tagText string
	row := db.QueryRow("SELECT tagid, value FROM tags WHERE taskid = ? LIMIT 1", 3)
	err = row.Scan(&tagId, &tagText)
	if err != nil {
		t.Fatal(err)
	}
	if tagId != 4 {
		t.Errorf("Error: got %v, want %v", tagId, 4)
	}
	if tagText != "Created tag1" {
		t.Errorf("Error: got %v, want %v", tagText, "Created tag1")
	}
}
