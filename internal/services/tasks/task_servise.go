package tasks

import (
	"fmt"
	e "go_rest/internal/errors"
	m "go_rest/internal/models"
	"go_rest/internal/repositories/tasks"

	"github.com/google/wire"
)

type Interface interface {
	Create(*m.Task) error
	Get(id int) (m.Task, error)
	Delete(id int) error
	DeleteAll() error
	GetByTag(tag string) (m.Tasks, error)
	GetByDate(date string) (m.Tasks, error)
	GetAll() (m.Tasks, error)
}

type TaskService struct {
	user       *m.User
	repository tasks.TaskRepositoryInterface
}

func (t *TaskService) SetUser(user *m.User) {
	t.user = user
}

func (t *TaskService) Create(task *m.Task) error {
	return t.repository.Create(task, t.user.Id)
}

func (t *TaskService) Get(id int) (m.Task, error) {
	if t.user.IsAdmin {
		return t.repository.Get(id)
	} else {
		return m.Task{}, fmt.Errorf("not implement err")
	}
}

func (t *TaskService) Delete(id int) error {
	if t.user.IsAdmin {
		return t.repository.Delete(id)
	} else {
		return fmt.Errorf("not implement err")
	}
}

func (t *TaskService) DeleteAll() error {
	if t.user.IsAdmin {
		return t.repository.DeleteAll()
	} else {
		return &e.UserNotAdminErr{UserId: t.user.Id}
	}

}

func (t *TaskService) GetByTag(tag string) (m.Tasks, error) {
	if t.user.IsAdmin {
		return t.repository.GetByTag(tag)
	} else {
		return m.Tasks{}, fmt.Errorf("not implement err")
	}
}

func (t *TaskService) GetByDate(date string) (m.Tasks, error) {
	if t.user.IsAdmin {
		return t.repository.GetByDate(date)
	} else {
		return m.Tasks{}, fmt.Errorf("not implement err")
	}
}

func (t *TaskService) GetAll() (m.Tasks, error) {
	if t.user.IsAdmin {
		return t.repository.GetAll()
	} else {
		return m.Tasks{}, fmt.Errorf("not implement err")
	}
}

func NewTaskService(r tasks.TaskRepositoryInterface) *TaskService {
	return &TaskService{
		repository: r,
	}
}

var TaskSet = wire.NewSet(
	tasks.NewTaskRepository,
	wire.Bind(new(tasks.TaskRepositoryInterface), new(*tasks.TaskRepository)),
	NewTaskService,
	wire.Bind(new(Interface), new(*TaskService)),
)
