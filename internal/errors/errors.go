package core

import (
	"errors"
	"fmt"
)

type ErrorWrapper struct {
	Message string
	Context string
	Err     error
}

func (err ErrorWrapper) Error() string {
	if err.Err != nil {
		return err.Err.Error()
	}
	return err.Message
}

func (err ErrorWrapper) Unwrap() error {
	return err.Err // Returns inner error
}

func (err ErrorWrapper) Dig() ErrorWrapper {
	var ew ErrorWrapper
	if errors.As(err.Err, &ew) {
		return ew.Dig()
	}
	return err
}

func NewErrorWrapper(context string, err error, message string) error {
	return ErrorWrapper{
		Message: message,
		Context: context,
		Err:     err,
	}
}

type InvalidPass struct {
	Email string
}

func (e *InvalidPass) Error() string {
	return fmt.Sprintf("User email %v not authorized. Invalid pass", e.Email)
}

func (e *InvalidPass) Is(target error) bool {
	t, ok := target.(*InvalidPass)
	if !ok {
		return false
	}
	return e.Email == t.Email
}

type UserNotAdminErr struct {
	UserId int
}

func (e *UserNotAdminErr) Error() string {
	return fmt.Sprintf("User id %v not admin", e.UserId)
}

func (e *UserNotAdminErr) Is(target error) bool {
	t, ok := target.(*UserNotAdminErr)
	if !ok {
		return false
	}
	return e.UserId == t.UserId
}

type TaskNotFound struct {
	TaskId int
}

func (e *TaskNotFound) Error() string {
	return fmt.Sprintf("Task id %v not found", e.TaskId)
}

func (e *TaskNotFound) Is(target error) bool {
	t, ok := target.(*TaskNotFound)
	if !ok {
		return false
	}
	return e.TaskId == t.TaskId
}

type EmptyTasks struct{}

func (e *EmptyTasks) Error() string {
	return "Empty Tasks"
}

func (e *EmptyTasks) Is(target error) bool {
	_, ok := target.(*EmptyTasks)
	return ok
}

type UserNotFound struct {
	UserId int
	Email  string
}

func (e *UserNotFound) Error() string {
	return fmt.Sprintf("User id %v, email %v not found", e.UserId, e.Email)
}

func (e *UserNotFound) Is(target error) bool {
	t, ok := target.(*UserNotFound)
	if !ok {
		return false
	}
	return e.UserId == t.UserId
}
