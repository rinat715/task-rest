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

// func (ts *TasksStoreSqlite) GetByTag(tag string) (m.Tasks, error) {
// 	var tasks m.Tasks

// 	stmt, err := ts.db.Prepare(`SELECT task.taskid, task.text, task.date, task.done
// 								FROM task INNER JOIN tag
// 								WHERE tag.value = ?`)
// 	if err != nil {
// 		return tasks, err
// 	}

// }

// func (ts *TasksStoreSqlite) GetByTag(tag string) (m.Tasks, error) {
// 	var tasks m.Tasks
// 	stmt, err := ts.db.Prepare(`SELECT task.taskid, task.taskid, task.date, task.done
// 							FROM task INNER JOIN tag
// 							WHERE task.taskid = tag.taskid`)
// 	if err != nil {
// 		return tasks, err
// 	}

// 	rows, err := stmt.Query(1)
// 	if err != nil {
// 		return tasks, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var task m.Task
//         if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist,
//             &alb.Price, &alb.Quantity); err != nil {
//             return albums, err
//         }
//         albums = append(albums, album)
// 		// ...
// 	}
// 	if err = rows.Err(); err != nil {
// 		return tasks, err
// 	}

// 	defer stmt.Close()
// 	return tasks, err

// }
// func (ts *TasksStoreSqlite) GetByDate(date m.JsonDate) m.Tasks {}

func (ts *tasks) getTasks() (m.Tasks, error) {
	var err error
	var tasks m.Tasks
	stmt, err := ts.db.Prepare("SELECT taskid, text, date, done FROM task")
	if err != nil {
		return tasks, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return tasks, err
	}

	defer stmt.Close()
	return parseTasks(rows)
}

func (ts *tasks) GetAll() (m.Tasks, error) {
	var err error
	var tasks m.Tasks
	tasks, err = ts.getTasks()
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
	return tasks, nil
}
