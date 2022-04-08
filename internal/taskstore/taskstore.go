// package taskstore provides a simple in-memory "data store" for tasks.
// Tasks are uniquely identified by numeric IDs.
//
// Eli Bendersky [https://eli.thegreenplace.net]
// This code is in the public domain.
package taskstore

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type Task struct {
	Id   int       `json:"id"`
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}

// десериализация
func (c Task) Serialize(data []byte) ([]Task, error) {
	tasks := []Task{}
	err := json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

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

	tasks  map[int]Task
	nextId int
}

// какой то аналог python dict values
func (ts *TaskStore) values() []Task {
	arr := make([]Task, 0, len(ts.tasks))
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
	ts := &TaskStore{}
	ts.tasks = make(map[int]Task) // создание среза
	ts.nextId = 0
	return ts
}

// CreateTask creates a new task in the store.
func (ts *TaskStore) CreateTask(text string, tags []string, due time.Time) int {
	//
	ts.Lock()         // лок стора  sync.Mutex
	defer ts.Unlock() // после выполнения фукции инлок

	task := Task{
		Id:   ts.nextId,
		Text: text,
		Due:  due}
	task.Tags = make([]string, len(tags)) // срез длиной количество тегов
	copy(task.Tags, tags)

	ts.tasks[ts.nextId] = task
	ts.nextId++
	return task.Id
}

// GetTask retrieves a task from the store, by id. If no such id exists, an
// error is returned.
func (ts *TaskStore) GetTask(id int) (Task, error) {
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
func (ts *TaskStore) DeleteTask(id int) error {
	ts.Lock()
	defer ts.Unlock()

	if _, ok := ts.tasks[id]; !ok {
		return fmt.Errorf("task with id=%d not found", id)
	}

	delete(ts.tasks, id)
	return nil
}

// DeleteAllTasks deletes all tasks in the store.
func (ts *TaskStore) DeleteAllTasks() error {
	ts.Lock()
	defer ts.Unlock()

	ts.tasks = make(map[int]Task)
	return nil
}

// GetAllTasks returns all the tasks in the store, in arbitrary order.
func (ts *TaskStore) GetAllTasks() []Task {
	ts.Lock()
	defer ts.Unlock()

	/*
		возвращаю копию массива потому что ?????
	*/
	// for _, task := range ts.tasks {
	// 	allTasks = append(allTasks, task)
	// }
	return ts.values()
}

// GetTasksByTag returns all the tasks that have the given tag, in arbitrary
// order.
func (ts *TaskStore) GetTasksByTag(tag string) []Task {
	ts.Lock()
	defer ts.Unlock()

	var tasks []Task

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
func (ts *TaskStore) GetTasksByDueDate(year int, month time.Month, day int) []Task {
	ts.Lock()
	defer ts.Unlock()

	var tasks []Task

	for _, task := range ts.tasks {
		y, m, d := task.Due.Date()
		if y == year && m == month && d == day {
			tasks = append(tasks, task)
		}
	}

	return tasks
}
