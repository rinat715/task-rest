// package taskstore provides a simple in-memory "data store" for tasks.
// Tasks are uniquely identified by numeric IDs.
//
// Eli Bendersky [https://eli.thegreenplace.net]
// This code is in the public domain.
package store

import (
	"fmt"
	m "go_rest/internal/models"
	"sync"
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

type TaskStore struct {
	sync.Mutex

	tasks  map[int]m.Task
	nextId int
}

// какой то аналог python dict values
func (ts *TaskStore) values() m.Tasks {
	arr := make(m.Tasks, 0, len(ts.tasks))
	for _, task := range ts.tasks {
		arr = append(arr, task)
	}
	return arr
}

// какой то аналог python dict[key] = value
func (ts *TaskStore) set(task m.Task) int {
	ts.tasks[ts.nextId] = task
	ts.nextId++
	return task.Id
}

func New() *TaskStore {
	return &TaskStore{
		tasks:  make(map[int]m.Task),
		nextId: 1,
	}
}

// CreateTask creates a new task in the store.
func (ts *TaskStore) Create(task *m.Task) {
	//
	ts.Lock()         // лок стора  sync.Mutex
	defer ts.Unlock() // после выполнения фукции инлок
	task.Id = ts.nextId
	ts.tasks[ts.nextId] = *task
	ts.nextId++
}

// GetTask retrieves a task from the store, by id. If no such id exists, an
// error is returned.
func (ts *TaskStore) Get(id int) (m.Task, error) {
	ts.Lock()
	defer ts.Unlock()

	t, ok := ts.tasks[id]
	if ok {
		return t, nil
	} else {
		return m.Task{}, fmt.Errorf("task with id=%d not found", id)
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

	ts.tasks = make(map[int]m.Task)
	return nil
}

// GetAllTasks returns all the tasks in the store, in arbitrary order.
func (ts *TaskStore) GetAll() m.Tasks {
	ts.Lock()
	defer ts.Unlock()
	return ts.values()
}

// GetTasksByTag returns all the tasks that have the given tag, in arbitrary
// order.
func (ts *TaskStore) GetByTag(tag string) m.Tasks {
	ts.Lock()
	defer ts.Unlock()
	var tasks m.Tasks

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
func (ts *TaskStore) GetByDate(date m.JsonDate) m.Tasks {
	ts.Lock()
	defer ts.Unlock()
	var tasks m.Tasks
	for _, task := range ts.tasks {
		if task.Date.Equal(date) {
			tasks = append(tasks, task)
		}
	}
	return tasks
}
