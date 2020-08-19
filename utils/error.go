package utils

// InternalError is when the proxy fails not the request
type InternalError struct {
	err     error
	Message string `json:"message"`
	Code    string `json:"code"`
}

// Error returns the error message
func (e *InternalError) Error() string {
	return e.Message
}

// Root returns the root error, if any
func (e *InternalError) Root() error {
	return e.err
}

// New creates a new internal error from an error
func New(err error, code, message string) error {
	if err != nil {
		message = err.Error()
	}
	return &InternalError{
		err:     err,
		Message: message,
		Code:    code,
	}
}
