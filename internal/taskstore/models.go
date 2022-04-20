package taskstore

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const JsonDateForm = "2006-01-02"

type JsonDate time.Time

func (t *JsonDate) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	nt, err := time.Parse(JsonDateForm, s)
	if err != nil {
		*t = JsonDate(time.Time{}) // в дату записывается начальное время
		return
	}
	*t = JsonDate(nt)
	return
}

func (t JsonDate) MarshalJSON() ([]byte, error) {
	return []byte(t.String()), nil
}

// String returns the time in the custom format
func (t JsonDate) String() string {
	nt := time.Time(t)
	return fmt.Sprintf("%q", nt.Format(JsonDateForm))
}

func unmarshal(data []byte, v interface{}) error {
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	return nil
}

func marshal(v interface{}) ([]byte, error) {
	js, err := json.Marshal(v)
	if err != nil {
		return []byte{}, err
	}
	return js, nil
}

type Serializer interface {
	Serialize() ([]byte, error)
}

type Task struct {
	Id   int      `json:"id,omitempty" default:"0"`
	Text string   `json:"text"`
	Tags []string `json:"tags"`
	Date JsonDate `json:"date"`
	Done bool     `json:"done"`
}

// десериализация
func (c Task) Deserialize(data []byte) (Task, error) {
	task := Task{}
	return task, unmarshal(data, &task)
}

func (c *Task) Serialize() ([]byte, error) {
	return marshal(c)
}

type Tasks []Task

// десериализация
func (c Tasks) Deserialize(data []byte) (Tasks, error) {
	tasks := []Task{}
	return tasks, unmarshal(data, &tasks)
}

func (c *Tasks) Serialize() ([]byte, error) {
	return marshal(c)
}
