package sqlstore

import (
	"database/sql"
	m "go_rest/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

func parseTasks(rows *sql.Rows) (m.Tasks, error) {
	var err error
	var tasks m.Tasks

	for rows.Next() {
		var task m.Task
		err := rows.Scan(&task.Id, &task.Text, &task.Date, &task.Done)
		tasks = append(tasks, task)
		if err != nil {
			return tasks, err
		}
	}
	err = rows.Err()
	if err != nil {
		return tasks, err
	}
	defer rows.Close()
	return tasks, err
}

func parseTagToArray(rows *sql.Rows) ([]string, error) {
	var err error
	var tags []string
	// я так понимаю заместо этого используют sqlx
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		tags = append(tags, tag)
		if err != nil {
			return tags, err
		}
	}
	err = rows.Err()
	if err != nil {
		return tags, err
	}
	defer rows.Close()
	return tags, err

}

func parseTagToMap(rows *sql.Rows) (map[int][]string, error) {
	var err error
	tags := make(map[int][]string)

	for rows.Next() {
		var tag string
		var taskid int
		err := rows.Scan(&taskid, &tag)
		i, ok := tags[taskid]
		if ok != true {
			tags[taskid] = []string{tag}
		} else {
			i = append(i, tag)
		}
		if err != nil {
			return tags, err
		}
	}
	err = rows.Err()
	if err != nil {
		return tags, err
	}
	defer rows.Close()
	return tags, err
}
