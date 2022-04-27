package sqlstore

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type tag struct {
	db *sql.DB
}

func (ts *tag) addTags(tx *sql.Tx, tags []string, taskid int64) error {
	stmt, err := tx.Prepare("INSERT INTO tag(value, taskid) VALUES(?, ?)")
	if err != nil {
		return err
	}
	for _, tag := range tags {
		_, err := stmt.Exec(tag, taskid)
		if err != nil {
			return err
		}
	}
	defer stmt.Close()
	return nil
}

func (ts *tag) getTagsbyTaskid(taskid int) ([]string, error) {

	stmt, err := ts.db.Prepare(`SELECT value FROM tag WHERE taskid = ?`)
	if err != nil {
		return []string{}, err
	}

	rows, err := stmt.Query(taskid)
	if err != nil {
		return []string{}, err
	}

	defer stmt.Close()
	return parseTagToArray(rows)
}

func (ts *tag) getTags() (map[int][]string, error) {
	var tags map[int][]string

	stmt, err := ts.db.Prepare(`SELECT taskid, value FROM tag`)
	if err != nil {
		return tags, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return tags, err
	}

	defer stmt.Close()
	return parseTagToMap(rows)
}
