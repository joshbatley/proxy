package fail

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Code error code type
type code error

// Error is when the proxy fails not the request
type Error struct {
	Inner   error
	Message string
	Code    code
}

var (
	// ErrMissingCol collection ID no in DB
	ErrMissingCol code = errors.New("collection_missing")
	// ErrURLInvalid requested URL is not valid
	ErrURLInvalid code = errors.New("URL_invalid")
	// ErrInternal internal error
	ErrInternal code = errors.New("internal_error")
	// ErrNoData when no data is found
	ErrNoData code = errors.New("no_data")
	// ErrOffline when no internet connection found
	ErrOffline code = errors.New("no_connection_found")
	// ErrResponse when no response found
	ErrResponse code = errors.New("no_response_found")
)

func (e *Error) Error() string {
	return fmt.Sprint(e.Message, " - ", e.Inner.Error())
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
			"error":   e.Inner.Error(),
		},
	)
}

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

// InternalError returns a unexpected error
func InternalError(err error) error {
	return &Error{
		Inner:   err,
		Code:    ErrInternal,
		Message: "Internal Error",
	}
}

// OfflineError returns a unexpected error
func OfflineError(err error) error {
	return &Error{
		Inner:   err,
		Code:    ErrOffline,
		Message: "Check you internet connection",
	}
}

// ResponseMissing returns a unexpected error
func ResponseMissing(err error) error {
	return &Error{
		Inner:   err,
		Code:    ErrResponse,
		Message: "Response not found, check preferred status type on endpoint",
	}
}
