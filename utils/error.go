package utils

import (
	"encoding/json"
	"errors"
)

// Code error code type
type code error

// InternalError is when the proxy fails not the request
type InternalError struct {
	Inner   error
	Message string
	Code    code
}

var (
	// ErrCollectionMissing collection ID no in DB
	ErrCollectionMissing code = errors.New("collection_missing")
)

// Error returns the error message
func (e *InternalError) Error() string {
	return e.Message
}

func (e *InternalError) Unwrap() error {
	return e.Code
}

// MarshalJSON Parse struct to json
func (e *InternalError) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		map[string]string{
			"message": e.Message,
			"code":    e.Code.Error(),
		},
	)
}

// New creates a new internal error from an error
func New(err error, c code, message string) error {
	return &InternalError{
		Inner:   err,
		Message: message,
		Code:    c,
	}
}

// ColMissingErr returns new colleciton missing error
func ColMissingErr(err error) error {
	return New(
		err,
		ErrCollectionMissing,
		"No collection by that ID exists",
	)
}
