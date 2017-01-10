package exception

import "fmt"

type Handler interface {
	Raise(err error, code int)
	Print() string
}

type ServerError struct {
	Error   error
	Message string
	Code    int
}

func (se ServerError) Raise(err error, code int) {
	se.Error = err
	se.Message = err.Error()
	se.Code = code
}

// Print implements the exception handler interface
func (se ServerError) Print() string {
	return fmt.Sprintf("[SERVER_ERROR][message: %s][code: %d]", se.Message, se.Code)
}
