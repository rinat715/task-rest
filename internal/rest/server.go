package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"strconv"

	"github.com/bmizerany/pat"

	"go_rest/internal/models"
)

// как бы класс только структура которая создает экземпляр сервера
// store это экземпляр хранилища тасков
type taskServer struct {
	store  models.Repository
	router *pat.PatternServeMux
}

func NewTaskServer() {
	s := &taskServer{
		store:  taskstore.New(),
		router: pat.New(),
	}
	s.routers()
	log.Fatal(http.ListenAndServe("localhost:"+os.Getenv("SERVERPORT"), s.router)) // пробрасываю порт который будет слушаать сервер
}

// utils
// рендрер json responce
// по факту это надо засунуть в миддлеваре
func (s *taskServer) jsonResponse(w http.ResponseWriter, v models.Serializer) error {
	js, err := v.Serialize()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return nil
}

func (s *taskServer) parseJsonRequest(w http.ResponseWriter, req *http.Request, v interface{}) error {
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return err
	}
	return nil
}

func (s *taskServer) validateRequestType(w http.ResponseWriter, req *http.Request, request_type string) error {
	// Enforce a JSON Content-Type.
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	if mediatype != request_type {
		http.Error(w, fmt.Sprint("expect %s Content-Type", request_type), http.StatusUnsupportedMediaType)
		return err
	}
	return nil
}

func (s *taskServer) getIdfromQuery(req *http.Request) int {
	param := req.URL.Query().Get(":taskid")
	taskid, err := strconv.Atoi(param)
	if err != nil {
		return 0
	}
	return taskid
}

// routes
func (s *taskServer) routers() {
	s.router.Get("/tasks", http.HandlerFunc(s.GetTasks))
	s.router.Get("/tasks/:taskid", http.HandlerFunc(s.GetTaskbyId))
	s.router.Post("/tasks", http.HandlerFunc(s.PostTask))
	s.router.Del("/tasks", http.HandlerFunc(s.DelTasks))
	s.router.Del("/tasks/:taskid", http.HandlerFunc(s.DelTaskbyId))
}

// handlers
func (s *taskServer) GetTasks(w http.ResponseWriter, req *http.Request) {
	var tasks models.Tasks
	tag := req.URL.Query().Get("tag")
	date := req.URL.Query().Get("date")

	if tag != "" && date != "" {
		http.Error(w, "Query params invalid", http.StatusBadRequest)
		return
	}

	if tag != "" {
		tasks = s.store.GetByTag(tag)
	}

	if date != "" {
		t, err := models.JsonDateParse(date)
		if err != nil {
			http.Error(w, fmt.Sprint("Date %s invalid", date), http.StatusBadRequest)
			return
		}
		tasks = s.store.GetByDate(t)
	}
	tasks = s.store.GetAll()
	s.jsonResponse(w, &tasks)
}

func (s *taskServer) GetTaskbyId(w http.ResponseWriter, req *http.Request) {
	var task models.Task
	taskid := s.getIdfromQuery(req)
	if taskid == 0 {
		http.Error(w, "Query params invalid", http.StatusBadRequest)
		return
	}
	task, err := s.store.Get(taskid)
	if err != nil {
		http.Error(w, fmt.Sprint("Task by id %s not found"), http.StatusNotFound)
	}
	s.jsonResponse(w, &task)

}

func (s *taskServer) PostTask(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling task create at %s\n", req.URL.Path)
	var err error

	// Types used internally in this handler to (de-)serialize the request and
	// response from/to JSON.
	err = s.validateRequestType(w, req, "application/json")
	if err != nil {
		return
	}

	var task models.Task
	err = s.parseJsonRequest(w, req, &task)
	if err != nil {
		return
	}
	s.store.Create(&task)
	w.WriteHeader(http.StatusCreated)
	s.jsonResponse(w, &task)

}

func (s *taskServer) DelTaskbyId(w http.ResponseWriter, req *http.Request) {
	taskid := s.getIdfromQuery(req)

	if taskid == 0 {
		http.Error(w, "Query params invalid", http.StatusBadRequest)
		return
	}

	err := s.store.Delete(taskid)
	if err != nil {
		http.Error(w, fmt.Sprint("Task by id %s not found"), http.StatusNotFound)
	}
}

func (s *taskServer) DelTasks(w http.ResponseWriter, req *http.Request) {
	s.store.DeleteAll()
}

func main() {
	log.Printf("Сервер запущен")
	NewTaskServer()
}