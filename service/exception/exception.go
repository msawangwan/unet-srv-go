package exception

import "fmt"

type Handler interface {
	Print(err error, code int) string
}

type ServerError struct {
	Error   error
	Message string
	Code    int
}

// Print implements the exception handler interface
func (se ServerError) Print(err error, code int) string {
	return fmt.Sprintf("[SERVER_ERROR][message: %s][code: %d]", err.Error(), code)
}
