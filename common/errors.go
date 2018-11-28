package common

import "fmt"

func NewTimeoutError(message string, args ...interface{}) TimeoutError {
	return TimeoutError{
		Message: fmt.Sprintf(message, args...),
	}
}

type TimeoutError struct {
	Message string
}

func (instance TimeoutError) Error() string {
	return instance.Message
}

func (instance TimeoutError) String() string {
	return instance.Error()
}

func IsTimeout(candidate error) bool {
	if _, ok := candidate.(*TimeoutError); ok {
		return true
	}
	if _, ok := candidate.(TimeoutError); ok {
		return true
	}
	return false
}
