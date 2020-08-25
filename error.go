package proxy

import (
	"encoding/json"
	"errors"
)

// Code error code type
type code error

// Error is when the proxy fails not the request
type Error struct {
	Inner   error
	Message string
	Code    code
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Code
}

// MarshalJSON Parse struct to json
func (e *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		map[string]string{
			"message": e.Message,
			"code":    e.Code.Error(),
		},
	)
}

var (
	// ErrMissingCol collection ID no in DB
	ErrMissingCol code = errors.New("collection_missing")
	// ErrURLInvalid requested URL is not valid
	ErrURLInvalid code = errors.New("URL_invalid")
)

// MissingColErr returns new colleciton missing error
func MissingColErr(err error) error {
	return &Error{
		Inner:   err,
		Code:    ErrMissingCol,
		Message: "No collection by that ID exists",
	}
}

// URLInvalidErr returns new colleciton missing error
func URLInvalidErr(err error) error {
	return &Error{
		Inner:   err,
		Code:    ErrURLInvalid,
		Message: "Requested URL is not valid",
	}
}
