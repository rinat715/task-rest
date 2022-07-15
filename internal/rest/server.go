package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"

	"go_rest/internal/config"
	"go_rest/internal/logger"
	"go_rest/internal/models"
	taskService "go_rest/internal/services/tasks"
	userService "go_rest/internal/services/users"
)

const QueryDateForm = "2006-01-02"

func NewHttpServer(s *taskServer) *http.Server {
	return &http.Server{
		Addr:    s.BuildUrl(),
		Handler: s.Router,
	}
}

// как бы класс только структура которая создает экземпляр сервера
// store это экземпляр хранилища тасков
type taskServer struct {
	UserService userService.Interface
	TaskService taskService.Interface
	Router      *pat.PatternServeMux
	*config.Config
}

func NewTaskServer(u userService.Interface, t taskService.Interface, r *pat.PatternServeMux, c *config.Config) *taskServer {
	s := &taskServer{
		UserService: u,
		TaskService: t,
		Router:      r,
		Config:      c,
	}
	s.routers()
	return s

}

// utils
func (s *taskServer) BuildUrl() string {
	return fmt.Sprintf("%v:%v", s.Host, s.Port)
}

// рендрер json responce
func (s *taskServer) jsonResponse(w http.ResponseWriter, v json.Marshaler) error {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return nil
}

func (s *taskServer) parseJsonRequest(w http.ResponseWriter, req *http.Request, v json.Unmarshaler) error {
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return err
	}
	return nil
}

func (s *taskServer) getIdfromQuery(req *http.Request) (int, error) {
	param := req.URL.Query().Get(":taskid")
	taskid, err := strconv.Atoi(param)
	if err != nil {
		return 0, fmt.Errorf("invalid taskid: %v", param)
	}
	return taskid, nil
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
		logger.Error(fmt.Sprintf("tag: %v, date: %v", tag, date))
		http.Error(w, "Query params invalid", http.StatusBadRequest)
		return
	}

	if tag != "" {
		tasks, err = s.TaskService.GetByTag(tag)
		if err != nil {
			logger.Error(err)
			http.Error(w, fmt.Sprintf("Store error:  %v", err), http.StatusServiceUnavailable)
			return
		}
	}

	if date != "" {
		t, err := time.Parse(QueryDateForm, date)
		if err != nil {
			logger.Error(err)
			http.Error(w, fmt.Sprintf("Date %v invalid", date), http.StatusBadRequest)
			return
		}
		tasks, err = s.TaskService.GetByDate(t.Format(QueryDateForm))
		if err != nil {
			logger.Error(err)
			http.Error(w, fmt.Sprintf("Store error:  %v", err), http.StatusServiceUnavailable)
			return
		}
	}
	tasks, err = s.TaskService.GetAll()
	if err != nil {
		logger.Error(err)
		http.Error(w, fmt.Sprintf("Store error:  %v", err), http.StatusServiceUnavailable)
		return
	}
	s.jsonResponse(w, &tasks)
}

func (s *taskServer) GetTaskbyId(w http.ResponseWriter, req *http.Request) {
	var task models.Task
	var err error
	taskid, err := s.getIdfromQuery(req)
	if err != nil {
		logger.Error(err)
		http.Error(w, "Query params invalid", http.StatusBadRequest)
		return
	}
	task, err = s.TaskService.Get(taskid)
	if err != nil {
		logger.Error(err)
		http.Error(w, fmt.Sprintf("Task by id %v not found", task.Id), http.StatusNotFound)
		return
	}
	s.jsonResponse(w, &task)

}

func (s *taskServer) PostTask(w http.ResponseWriter, req *http.Request) {
	var err error
	var task models.Task

	err = s.parseJsonRequest(w, req, &task)
	if err != nil {
		logger.Error(err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	err = s.TaskService.Create(&task)
	if err != nil {
		logger.Error(err)
		http.Error(w, "Failed create task", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	s.jsonResponse(w, &task)

}

func (s *taskServer) DelTaskbyId(w http.ResponseWriter, req *http.Request) {
	var err error
	taskid, err := s.getIdfromQuery(req)

	if err != nil {
		logger.Error(err)
		http.Error(w, "Query params invalid", http.StatusBadRequest)
		return
	}

	err = s.TaskService.Delete(taskid)
	if err != nil {
		logger.Error(err)
		http.Error(w, fmt.Sprintf("Task by id %v not found", taskid), http.StatusNotFound)
		return
	}
}

func (s *taskServer) DelTasks(w http.ResponseWriter, req *http.Request) {
	err := s.TaskService.DeleteAll()
	if err != nil {
		logger.Error(err)
		http.Error(w, "Failed delete all tasks", http.StatusBadRequest)
		return
	}

}
