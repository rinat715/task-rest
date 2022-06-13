package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"go_rest/internal/taskstore/sqlstore"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"

	"go_rest/internal/logger"
	"go_rest/internal/models"
)

// как бы класс только структура которая создает экземпляр сервера
// store это экземпляр хранилища тасков
type taskServer struct {
	store  models.Repository
	Router *pat.PatternServeMux
}

func NewTaskServer(ts *sqlstore.Store) *taskServer {
	s := &taskServer{
		store:  ts,
		Router: pat.New(),
	}
	s.routers()
	return s

}

// utils
// рендрер json responce
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
	commonHandlers := alice.New(PanicRecovery, Logging, ValidateRequestJsonType)
	s.Router.Get("/tasks", commonHandlers.Then(http.HandlerFunc(s.GetTasks)))
	s.Router.Get("/tasks/:taskid", commonHandlers.Then(http.HandlerFunc(s.GetTaskbyId)))
	s.Router.Post("/tasks", commonHandlers.Then(http.HandlerFunc(s.PostTask)))
	s.Router.Del("/tasks", commonHandlers.Then(http.HandlerFunc(s.DelTasks)))
	s.Router.Del("/tasks/:taskid", commonHandlers.Then(http.HandlerFunc(s.DelTaskbyId)))
	s.Router.Get("/test_panic", commonHandlers.Then(http.HandlerFunc(s.TestPanic)))
}

// handlers
func (s *taskServer) TestPanic(w http.ResponseWriter, req *http.Request) {
	panic("Error: Тестовый роут для проверки паники")
}

func (s *taskServer) GetTasks(w http.ResponseWriter, req *http.Request) {
	tag := req.URL.Query().Get("tag")
	date := req.URL.Query().Get("date")
	var tasks models.Tasks
	var err error

	if tag != "" && date != "" {
		http.Error(w, "Query params invalid", http.StatusBadRequest)
		return
	}

	if tag != "" {
		tasks, err = s.store.GetByTag(tag)
		if err != nil {
			http.Error(w, fmt.Sprintf("Store error:  %v", err), http.StatusServiceUnavailable)
			return
		}
	}

	if date != "" {
		t, err := models.JsonDateParse(date)
		if err != nil {
			http.Error(w, fmt.Sprintf("Date %v invalid", date), http.StatusBadRequest)
			return
		}
		tasks, err = s.store.GetByDate(t)
		if err != nil {
			http.Error(w, fmt.Sprintf("Store error:  %v", err), http.StatusServiceUnavailable)
			return
		}
	}
	tasks, err = s.store.GetAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("Store error:  %v", err), http.StatusServiceUnavailable)
		return
	}
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
		http.Error(w, fmt.Sprintf("Task by id %v not found", task.Id), http.StatusNotFound)
	}
	s.jsonResponse(w, &task)

}

func (s *taskServer) PostTask(w http.ResponseWriter, req *http.Request) {
	var err error
	var task models.Task

	err = s.parseJsonRequest(w, req, &task)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	err = s.store.Create(&task)
	if err != nil {
		logger.Error(err)
		http.Error(w, "Failed create task", http.StatusBadRequest)
		return
	}
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
		http.Error(w, fmt.Sprintf("Task by id %v not found", taskid), http.StatusNotFound)
	}
}

func (s *taskServer) DelTasks(w http.ResponseWriter, req *http.Request) {
	s.store.DeleteAll()
}
