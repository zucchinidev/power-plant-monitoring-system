package errors

import "github.com/pkg/errors"

type notFound struct {
	error
}

func WrapNotFound(err error, format string, args ...interface{}) error {
	return &notFound{errors.Wrapf(err, format, args...)}
}

func NewNotFound(format string, args ...interface{}) error {
	return &notFound{errors.Errorf(format, args...)}
}

func IsNotFoundError(err error) bool {
	err = errors.Cause(err)
	_, ok := err.(*notFound)
	return ok
}