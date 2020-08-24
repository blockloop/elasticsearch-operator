package errors

import (
	"errors"
	"fmt"
)

// New returns a new structured error
func New(message string, keyValuePairs ...interface{}) *StructuredError {
	return &StructuredError{
		msg:           message,
		keyValuePairs: keyValuePairs,
	}
}

// Wrap wraps an error as a structured error
func Wrap(err error, message string, keyValuePairs ...interface{}) *StructuredError {
	return &StructuredError{
		cause:         err,
		msg:           message,
		keyValuePairs: keyValuePairs,
	}
}

// Error provides information about StructuredError
type Error interface {
	KVs() []interface{}
	Cause() error
	Error() string
}

// KVs returns all of the key/value pairs associated with the provided error. This is used for logging
//
// If err is wrapped then this will traverse the entire chain of errors appending all key/value pairs
// until it reaches the end. Because key/value pairs can be duplicated, newer errors
// take precedence over older errors when log is used. For example:
//
// e1 := errors.New("e1", "name", "cats")
// e2 := errors.Wrap(e1, "e2", "name", "dogs")
//
// In this scenario the logged field would be "name=dogs" because it is the newest value.
func KVs(err error) []interface{} {
	if err == nil {
		return nil
	}

	current := err
	kvs := make([]interface{}, 0)
	for s, ok := current.(Error); ok; s, ok = current.(Error) {
		if !ok {
			return kvs
		}
		kvs = append(kvs, s.KVs()...)
		current = s.Cause()
	}
	return kvs
}

// StructuredError is a structured error
type StructuredError struct {
	msg           string
	cause         error
	keyValuePairs []interface{}
}

func (e StructuredError) KVs() []interface{} {
	return e.keyValuePairs
}

func (e StructuredError) Cause() error {
	return e.cause
}

// Error returns the message and cause as a string
func (e StructuredError) Error() string {
	if e.cause != nil {
		return fmt.Errorf("%s: %w", e.msg, e.cause).Error()
	}
	return e.msg
}

// String returns the message and cause as a string
func (e StructuredError) String() string {
	return e.Error()
}

func (e StructuredError) Unwrap() error {
	var current error = e
	for s, ok := current.(Error); ok; s, ok = current.(Error) {
		if !ok || s.Cause() == nil {
			return current
		}
		current = s.Cause()
	}
	return current
}

// Unwrap is a wrapper for stdlib errors.Unwrap.
// It traverses the entire chain of wrapped errors to give you the original cause.
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// Is is a wrapper for stdlib errors.Is
// This is a shortcut for e.g. Unwrap(err) == io.EOF
func Is(err, target error) bool {
	return errors.Is(err, target)
}
