package apperror

import "fmt"

// AppError defines a structured application error.
type AppError struct {
	Code    string // machine readable code
	Message string // human readable message
	Err     error  // underlying error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

func New(code, msg string, err error) *AppError {
	return &AppError{Code: code, Message: msg, Err: err}
}
