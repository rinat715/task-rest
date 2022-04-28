// по сути создаем аналог data store in-memory
package sqlstore

import (
	"database/sql"
	m "go_rest/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type tasks struct {
	db *sql.DB
	*tag
}

func (ts *tasks) Create(task *m.Task) error {
	var err error

	tx, err := ts.db.Begin() // начало транзакции
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO task(text, date, done) VALUES(?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}

	res, err := stmt.Exec(task.Text, task.Date, task.Done)
	if err != nil {
		tx.Rollback()
		return err
	}
	Taskid, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	err = ts.addTags(tx, task.Tags, Taskid)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit() // при успешном завершении
	task.Id = int(Taskid)
	defer stmt.Close()
	return err
}

func (ts *tasks) Get(id int) (m.Task, error) {
	var err error
	var task m.Task

	stmt, err := ts.db.Prepare(`SELECT taskid, text, date, done
								FROM task
								WHERE taskid = ?`)
	if err != nil {
		return task, err
	}

	err = stmt.QueryRow(id).Scan(&task.Id, &task.Text, &task.Date, &task.Done)
	if err != nil {
		return task, err
	}

	tags, err := ts.getTagsbyTaskid(task.Id)
	if err != nil {
		return task, err
	}
	task.Tags = tags
	defer stmt.Close()
	return task, err

}
func (ts *tasks) Delete(id int) error {
	var err error
	res, err := ts.db.Exec("DELETE FROM task WHERE taskid = ?", id)
	if err != nil {
		return err
	}
	Taskid, err := res.LastInsertId()
	if err != nil {
		return err
	}
	_, err = ts.db.Exec("DELETE FROM tag WHERE taskid = ?", Taskid)
	if err != nil {
		return err
	}
	return nil
}

func (ts *tasks) DeleteAll() error {
	var err error
	_, err = ts.db.Exec("DELETE FROM task")
	if err != nil {
		return err
	}
	_, err = ts.db.Exec("DELETE FROM tag")
	if err != nil {
		return err
	}
	return nil
}

func (ts *tasks) GetByTag(tag string) (m.Tasks, error) {
	const query string = `SELECT task.taskid, task.text, task.date, task.done
						FROM task INNER JOIN tag
						WHERE task.taskid = tag.taskid
						AND tag.value = ?`
	return ts.getTasks(query, tag)
}

func (ts *tasks) GetByDate(date m.JsonDate) (m.Tasks, error) {
	const query string = `SELECT taskid, text, date, done
						FROM task
						WHERE date = ?`

	return ts.getTasks(query, date)
}

func (ts *tasks) GetAll() (m.Tasks, error) {
	var err error
	var tasks m.Tasks
	const query string = "SELECT taskid, text, date, done FROM task"

	stmt, err := ts.db.Prepare(query)
	if err != nil {
		return tasks, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return tasks, err
	}
	tasks, err = parseTasks(rows)
	if err != nil {
		return tasks, err
	}

	tags, err := ts.getTags()
	if err != nil {
		return tasks, err
	}
	for _, task := range tasks {
		task.Tags = tags[task.Id]
	}

	defer stmt.Close()
	return tasks, nil
}

func (ts *tasks) getTasks(query string, arg ...interface{}) (m.Tasks, error) {
	var err error
	var tasks m.Tasks

	stmt, err := ts.db.Prepare(query)
	if err != nil {
		return tasks, err
	}
	rows, err := stmt.Query(arg...)
	if err != nil {
		return tasks, err
	}

	tasks, err = parseTasks(rows)
	if err != nil {
		return tasks, err
	}

	defer stmt.Close()

	tags, err := ts.getTagsbyTaskids(tasks.GetIds())
	if err != nil {
		return tasks, err
	}

	return tasks.AddTags(tags), nil
}
