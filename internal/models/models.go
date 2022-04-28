package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const JsonDateForm = "2006-01-02"

func JsonDateParse(value string) (JsonDate, error) {
	nt, err := time.Parse(JsonDateForm, value)
	if err != nil {
		return JsonDate(time.Time{}), err
	}
	return JsonDate(nt), nil
}

type JsonDate time.Time

func (t *JsonDate) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	*t, err = JsonDateParse(s)
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

func (t JsonDate) Equal(u JsonDate) bool {
	nt := time.Time(t)
	return nt.Equal(time.Time(u))
}

// это для того чтобы реализовать конвертирование JsonDate в sql
func (t JsonDate) Value() (driver.Value, error) {
	return time.Time(t), nil
}

// https://bun.uptrace.dev/guide/custom-types.html#sql-scanner
func (t *JsonDate) Scan(value interface{}) error {
	if value == nil {
		*t = JsonDate(time.Time{})
		return nil
	}
	res, ok := value.(time.Time)
	if ok == false {
		*t = JsonDate(time.Time{})
		return errors.New("Invalid type for JsonDate.Scan")
	}
	*t = JsonDate(res)
	return nil
}

type Tags []string

func (t *Tags) Scan(value interface{}) error {
	if value == nil {
		*t = Tags([]string{})
		return nil
	}
	res, ok := value.(Tags)
	if ok == false {
		*t = Tags([]string{})
		return errors.New("Invalid type for Tags.Scan")
	}
	*t = Tags(res)
	return nil
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
	Id   int      `json:"id" default:"0"`
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

func (c Tasks) GetIds() []int {
	var ids []int
	for _, i := range c {
		ids = append(ids, i.Id)
	}
	return ids
}

func (c Tasks) AddTags(tags map[int][]string) Tasks {
	var res Tasks
	for _, task := range c {
		task.Tags = tags[task.Id]
		res = append(res, task)
	}
	return res
}

// все репозитории соответствуют этому интерфейсу
type Repository interface {
	Create(*Task) error
	Get(id int) (Task, error)
	Delete(id int) error
	DeleteAll() error
	GetByTag(tag string) (Tasks, error)
	GetByDate(date JsonDate) (Tasks, error)
	GetAll() (Tasks, error)
}
