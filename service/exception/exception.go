package exception

import "fmt"

// Handler is an error type interface
type Handler interface {
	Raise(err error, code int)
	Print() string
}

// ServerError is a server error abstraction
type ServerError struct {
	Error   error
	Message string
	Code    int
}

// Raise implements the exception handler interface
func (se *ServerError) Raise(err error, code int) {
	se.Error = err
	se.Message = err.Error()
	se.Code = code
}

// Print implements the exception handler interface
func (se *ServerError) Print() string {
	return fmt.Sprintf("[SERVER_ERROR][message: %s][code: %d]", se.Message, se.Code)
}
