package core

import (
	"errors"
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
