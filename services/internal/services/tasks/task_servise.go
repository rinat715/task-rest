package tasks

import (
	m "go_rest/internal/models"
	"go_rest/internal/repositories/tasks"

	"github.com/google/wire"
)

type Interface interface {
	Create(task *m.Task, userId int) error
	Get(id int) (m.Task, error)
	Delete(id int) error
	DeleteAll() error
	GetByTag(tag string) (m.Tasks, error)
	GetByDate(date string) (m.Tasks, error)
	GetAll() (m.Tasks, error)
}

type TaskService struct {
	user       *m.User // TODO переделать
	repository tasks.TaskRepositoryInterface
}

func (t *TaskService) SetUser(user *m.User) {
	t.user = user
}

func (t *TaskService) Create(task *m.Task, userId int) error {
	return t.repository.Create(task, userId)
}

func (t *TaskService) Get(id int) (m.Task, error) {
	return t.repository.Get(id)
}

func (t *TaskService) Delete(id int) error {
	return t.repository.Delete(id)
}

func (t *TaskService) DeleteAll() error {
	return t.repository.DeleteAll()
}

func (t *TaskService) GetByTag(tag string) (m.Tasks, error) {
	return t.repository.GetByTag(tag)
}

func (t *TaskService) GetByDate(date string) (m.Tasks, error) {
	return t.repository.GetByDate(date)
}

func (t *TaskService) GetAll() (m.Tasks, error) {
	return t.repository.GetAll()
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
