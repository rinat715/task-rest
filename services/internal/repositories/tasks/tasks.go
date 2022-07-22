package tasks

import (
	"database/sql"
	m "go_rest/internal/models"
)

var baseQuery string = `SELECT taskid, userid, text, date, done, tags.tagid, tags.value
						FROM tasks
						LEFT JOIN tags USING(taskid)`

type TaskRepositoryInterface interface {
	Create(task *m.Task, userId int) error
	Delete(id int) error
	DeleteAll() error

	Get(userId int) (m.Task, error)

	GetByTag(tag string) (m.Tasks, error)
	GetByDate(date string) (m.Tasks, error)
	GetAll() (m.Tasks, error)
}

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db}
}

func (ts *TaskRepository) parseTask(rows *sql.Rows) (m.Task, error) {
	var err error
	var task m.Task
	var tags m.Tags

	for rows.Next() {
		var tag m.Tag
		err = rows.Scan(&task.Id, &task.UserId, &task.Text, &task.Date, &task.Done, &tag.Id, &tag.Text)
		if err != nil {
			return m.Task{}, err
		}
		tags = append(tags, tag)
	}
	err = rows.Err()
	if err != nil {
		return m.Task{}, err
	}
	defer rows.Close()
	task.Tags = tags
	return task, nil
}

func (ts *TaskRepository) parseTasks(rows *sql.Rows) (m.Tasks, error) {
	var err error
	var tasks m.Tasks
	var tasksMap = make(map[int]m.Task)

	for rows.Next() {
		var task m.Task
		var tagId sql.NullInt32
		var tagText sql.NullString

		err = rows.Scan(&task.Id, &task.UserId, &task.Text, &task.Date, &task.Done, &tagId, &tagText)
		if err != nil {
			return tasks, err
		}
		item, ok := tasksMap[task.Id]
		if ok {
			if tagId.Valid {
				item.Tags = append(item.Tags, m.Tag{
					Id:   int(tagId.Int32),
					Text: tagText.String,
				})
			}
		} else {
			tasksMap[task.Id] = task
			if tagId.Valid {
				item.Tags = append(item.Tags, m.Tag{
					Id:   int(tagId.Int32),
					Text: tagText.String,
				})
			}
		}
	}
	err = rows.Err()
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for _, task := range tasksMap {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (ts *TaskRepository) getTasks(query string, req string) (m.Tasks, error) {
	var tasks m.Tasks
	stmt, err := ts.db.Prepare(query)
	if err != nil {
		return tasks, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(req)
	if err != nil {
		return tasks, err
	}
	tasks, err = ts.parseTasks(rows)
	if err != nil {
		return tasks, err
	}
	return tasks, nil
}

func (ts *TaskRepository) GetAll() (m.Tasks, error) {
	var tasks m.Tasks

	rows, err := ts.db.Query(baseQuery)
	if err != nil {
		return tasks, err
	}
	return ts.parseTasks(rows)
}

func (ts *TaskRepository) GetByTag(tag string) (m.Tasks, error) {
	query := baseQuery + `WHERE tags.value = ?`
	return ts.getTasks(query, tag)
}

func (ts *TaskRepository) GetByDate(date string) (m.Tasks, error) {
	query := baseQuery + `WHERE tasks.date = ?`
	return ts.getTasks(query, date)
}

type transaction struct {
	tx *sql.Tx
}

func (t *transaction) exec(query string, args ...any) (int, error) {
	res, err := t.tx.Exec(query, args...)
	if err != nil {
		t.tx.Rollback()
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		t.tx.Rollback()
		return 0, err
	}
	return int(id), nil
}

func (t *transaction) commit() error {
	return t.tx.Commit()
}

func newTransaction(db *sql.DB) (*transaction, error) {
	if tx, err := db.Begin(); err != nil {
		return &transaction{}, err
	} else {
		return &transaction{tx}, err
	}
}

func (ts *TaskRepository) Get(taskid int) (m.Task, error) {
	var task m.Task
	query := baseQuery + `WHERE taskid = ?`

	rows, err := ts.db.Query(query, taskid)
	if err != nil {
		return task, err
	}
	task, err = ts.parseTask(rows)
	if err != nil {
		return task, err
	}
	return task, nil
}

func (ts *TaskRepository) Delete(id int) error {
	var err error
	taskQuery := "DELETE FROM tasks WHERE taskid = ?"
	tagQuery := "DELETE FROM tags WHERE taskid = ?"

	tx, err := newTransaction(ts.db)
	if err != nil {
		return err
	}

	_, err = tx.exec(taskQuery, id)
	if err != nil {
		return err
	}

	_, err = tx.exec(tagQuery, id)
	if err != nil {
		return err
	}

	err = tx.commit()
	if err != nil {
		return err
	}
	return nil
}

func (ts *TaskRepository) DeleteAll() error {
	var err error
	taskQuery := "DELETE FROM tasks"
	tagQuery := "DELETE FROM tags"

	tx, err := newTransaction(ts.db)
	if err != nil {
		return err
	}

	_, err = tx.exec(taskQuery)
	if err != nil {
		return err
	}
	_, err = tx.exec(tagQuery)
	if err != nil {
		return err
	}

	err = tx.commit()
	if err != nil {
		return err
	}
	return nil
}

func (ts *TaskRepository) Create(task *m.Task, userId int) error {
	var err error
	taskQuery := "INSERT INTO tasks(text, date, done, userid) VALUES(?, ?, ?, ?)"
	tagQuery := "INSERT INTO tags(value, taskid) VALUES(?, ?)"

	tx, err := ts.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(taskQuery)
	if err != nil {
		tx.Rollback()
		return err
	}
	res, err := stmt.Exec(task.Text, task.Date, task.Done, userId)
	if err != nil {
		tx.Rollback()
		return err
	}
	taskId, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}
	task.Id = int(taskId)
	stmt.Close()

	stmt, err = tx.Prepare(tagQuery)
	if err != nil {
		tx.Rollback()
		return err
	}
	for idx := range task.Tags {
		tag := &task.Tags[idx]
		res, err = stmt.Exec(tag.Text, task.Id)
		if err != nil {
			tx.Rollback()
			return err
		}
		tagId, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
			return err
		}
		tag.Id = int(tagId)
	}
	stmt.Close()

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
