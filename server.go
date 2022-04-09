package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"go_rest/internal/taskstore"
)

var urls  

type RequestTask struct {
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}

// рендрер json responce
func renderJsonResponse(w http.ResponseWriter, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func validateJsonRequestType(w http.ResponseWriter, req *http.Request) {
	// Enforce a JSON Content-Type.
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}
}

func parseJsonRequest(req *http.Request) (RequestTask, error) {
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var rt RequestTask
	if err := dec.Decode(&rt); err != nil {
		return rt, err
	}
	return rt, nil
}

// как бы класс только структура которая создает экземпляр сервера
// store это экземпляр хранилища тасков
type taskServer struct {
	store *taskstore.TaskStore
}

func NewTaskServer() *taskServer {
	store := taskstore.New() // создание
	return &taskServer{store: store}
}

// ну тут проверется path запроса
func (ts *taskServer) taskHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/task/" {
		// Request is plain "/task/", without trailing ID.
		if req.Method == http.MethodPost {
			ts.createTaskHandler(w, req)
		} else if req.Method == http.MethodGet {
			ts.getAllTasksHandler(w, req)
		} else if req.Method == http.MethodDelete {
			ts.deleteAllTasksHandler(w, req)
			// ну если метод не GET, DELETE or POST то ошибка метод не разрешен
		} else {
			http.Error(w, fmt.Sprintf("expect method GET, DELETE or POST at /task/, got %v", req.Method), http.StatusMethodNotAllowed)
			return
		}
	} else {
		// Request has an ID, as in "/task/<id>".
		path := strings.Trim(req.URL.Path, "/")
		pathParts := strings.Split(path, "/")
		if len(pathParts) < 2 {
			http.Error(w, "expect /task/<id> in task handler", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Method == http.MethodDelete {
			ts.deleteTaskHandler(w, req, int(id))
		} else if req.Method == http.MethodGet {
			ts.getTaskHandler(w, req, int(id))
		} else {
			http.Error(w, fmt.Sprintf("expect method GET or DELETE at /task/<id>, got %v", req.Method), http.StatusMethodNotAllowed)
			return
		}
	}
}

func (ts *taskServer) createTaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling task create at %s\n", req.URL.Path)

	// Types used internally in this handler to (de-)serialize the request and
	// response from/to JSON.

	validateJsonRequestType(w, req)

	type ResponseId struct {
		Id int `json:"id"`
	}
	rt, err := parseJsonRequest(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	id := ts.store.CreateTask(rt.Text, rt.Tags, rt.Due)
	var response = ResponseId{Id: id}
	renderJsonResponse(w, response)

}

func (ts *taskServer) getAllTasksHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get all tasks at %s\n", req.URL.Path)

	response := ts.store.GetAllTasks()
	renderJsonResponse(w, response)
}

func (ts *taskServer) getTaskHandler(w http.ResponseWriter, req *http.Request, id int) {
	log.Printf("handling get task at %s\n", req.URL.Path)

	task, err := ts.store.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJsonResponse(w, task)
}

func (ts *taskServer) deleteTaskHandler(w http.ResponseWriter, req *http.Request, id int) {
	log.Printf("handling delete task at %s\n", req.URL.Path)

	err := ts.store.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func (ts *taskServer) deleteAllTasksHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling delete all tasks at %s\n", req.URL.Path)
	ts.store.DeleteAllTasks()
}

func (ts *taskServer) tagHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling tasks by tag at %s\n", req.URL.Path)

	if req.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("expect method GET /tag/<tag>, got %v", req.Method), http.StatusMethodNotAllowed)
		return
	}

	path := strings.Trim(req.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 2 {
		http.Error(w, "expect /tag/<tag> path", http.StatusBadRequest)
		return
	}
	tag := pathParts[1]

	tasks := ts.store.GetTasksByTag(tag)
	renderJsonResponse(w, tasks)
}

func (ts *taskServer) dueHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling tasks by due at %s\n", req.URL.Path)

	if req.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("expect method GET /due/<date>, got %v", req.Method), http.StatusMethodNotAllowed)
		return
	}

	path := strings.Trim(req.URL.Path, "/")
	pathParts := strings.Split(path, "/")

	badRequestError := func() {
		http.Error(w, fmt.Sprintf("expect /due/<year>/<month>/<day>, got %v", req.URL.Path), http.StatusBadRequest)
	}
	if len(pathParts) != 4 {
		badRequestError()
		return
	}

	year, err := strconv.Atoi(pathParts[1])
	if err != nil {
		badRequestError()
		return
	}
	month, err := strconv.Atoi(pathParts[2])
	if err != nil || month < int(time.January) || month > int(time.December) {
		badRequestError()
		return
	}
	day, err := strconv.Atoi(pathParts[3])
	if err != nil {
		badRequestError()
		return
	}

	tasks := ts.store.GetTasksByDueDate(year, time.Month(month), day)
	renderJsonResponse(w, tasks)
}

func main() {
	mux := http.NewServeMux()
	server := NewTaskServer()
	mux.HandleFunc("/task/", server.taskHandler)
	mux.HandleFunc("/tag/", server.tagHandler)
	mux.HandleFunc("/due/", server.dueHandler)
	log.Printf("Сервер запущен")
	log.Fatal(http.ListenAndServe("localhost:"+os.Getenv("SERVERPORT"), mux)) // пробрасываю порт который будет слушаать сервер

}
