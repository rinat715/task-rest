// package taskstore provides a simple in-memory "data store" for tasks.
// Tasks are uniquely identified by numeric IDs.
//
// Eli Bendersky [https://eli.thegreenplace.net]
// This code is in the public domain.
package taskstore

import (
	"fmt"
	"sync"
	"time"
)

// TaskStore is a simple in-memory database of tasks; TaskStore methods are
// safe to call concurrently.
// sync.Mutex - мьютекс

/*
пример из руководства

type SafeCounter struct {
	mu sync.Mutex
	v  map[string]int
}

лок c.mu.Lock()
инлок c.mu.Unlock()
*/

// все репозитории соответствуют этому интерфейсу
type Repository interface {
	Create(text string, tags []string, date JsonDate, done bool) Task
	Get(id int) (Task, error)
	Delete(id int) error
	DeleteAll() error
	GetByTag(tag string) Tasks
	GetByDate(year int, month time.Month, day int) Tasks
	GetAll() Tasks
}

type TaskStore struct {
	sync.Mutex

	tasks  map[int]Task
	nextId int
}

// какой то аналог python dict values
func (ts *TaskStore) values() Tasks {
	arr := make(Tasks, 0, len(ts.tasks))
	for _, task := range ts.tasks {
		arr = append(arr, task)
	}
	return arr
}

// какой то аналог python dict[key] = value
func (ts *TaskStore) set(task Task) int {
	ts.tasks[ts.nextId] = task
	ts.nextId++
	return task.Id
}

func New() *TaskStore {
	return &TaskStore{
		tasks:  make(map[int]Task),
		nextId: 0,
	}
}

// CreateTask creates a new task in the store.
func (ts *TaskStore) Create(text string, tags []string, date JsonDate, done bool) Task {
	//
	ts.Lock()         // лок стора  sync.Mutex
	defer ts.Unlock() // после выполнения фукции инлок

	task := Task{
		Id:   ts.nextId,
		Text: text,
		Date: date,
		Done: done,
	}
	task.Tags = make([]string, len(tags)) // срез длиной количество тегов
	copy(task.Tags, tags)

	ts.tasks[ts.nextId] = task
	ts.nextId++
	return task
}

// GetTask retrieves a task from the store, by id. If no such id exists, an
// error is returned.
func (ts *TaskStore) Get(id int) (Task, error) {
	ts.Lock()
	defer ts.Unlock()

	t, ok := ts.tasks[id]
	if ok {
		return t, nil
	} else {
		return Task{}, fmt.Errorf("task with id=%d not found", id)
	}
}

// DeleteTask deletes the task with the given id. If no such id exists, an error
// is returned.
func (ts *TaskStore) Delete(id int) error {
	ts.Lock()
	defer ts.Unlock()

	if _, ok := ts.tasks[id]; !ok {
		return fmt.Errorf("task with id=%d not found", id)
	}

	delete(ts.tasks, id)
	return nil
}

// DeleteAllTasks deletes all tasks in the store.
func (ts *TaskStore) DeleteAll() error {
	ts.Lock()
	defer ts.Unlock()

	ts.tasks = make(map[int]Task)
	return nil
}

// GetAllTasks returns all the tasks in the store, in arbitrary order.
func (ts *TaskStore) GetAll() Tasks {
	ts.Lock()
	defer ts.Unlock()
	return ts.values()
}

// GetTasksByTag returns all the tasks that have the given tag, in arbitrary
// order.
func (ts *TaskStore) GetByTag(tag string) Tasks {
	ts.Lock()
	defer ts.Unlock()
	var tasks Tasks

taskloop:
	for _, task := range ts.tasks {
		for _, taskTag := range task.Tags {
			if taskTag == tag {
				tasks = append(tasks, task)
				continue taskloop
			}
		}
	}
	return tasks
}

// GetTasksByDueDate returns all the tasks that have the given due date, in
// arbitrary order.
func (ts *TaskStore) GetByDate(year int, month time.Month, day int) Tasks {
	ts.Lock()
	defer ts.Unlock()

	var tasks Tasks

	for _, task := range ts.tasks {
		nt := time.Time(task.Date)
		y, m, d := nt.Date()
		if y == year && m == month && d == day {
			tasks = append(tasks, task)
		}
	}

	return tasks
}
