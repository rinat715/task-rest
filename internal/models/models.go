package models

import (
	"encoding/json"
	"time"
)

type User struct {
	Id      int    `json:"id" default:"0"`
	Email   string `json:"email"`
	Pass    string `json:"-"`
	IsAdmin bool   `json:"is_admin" default:"false"`
}

type Users []User

func (v Users) MarshalJSON() ([]byte, error) {
	return json.Marshal(v)
}

type Task struct {
	Id     int       `json:"id" default:"0"`
	UserId int       `json:"-" default:"0"`
	Text   string    `json:"text"`
	Tags   Tags      `json:"tags"`
	Date   time.Time `json:"date"`
	Done   bool      `json:"done"`
}

type Tasks []Task

func (v Tasks) MarshalJSON() ([]byte, error) {
	return json.Marshal(v)
}

type Tag struct {
	Id   int    `json:"id" default:"0"`
	Text string `json:"text"`
}

type Tags []Tag

func (v Tags) MarshalJSON() ([]byte, error) {
	return json.Marshal(v)
}
