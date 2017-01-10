package handler

import (
	"errors"

	"net/http"

	"github.com/msawangwan/unet-srv-go/service/exception"
)

var (
	errNilBody = errors.New("handler: the request body was nil")
)

func checkBody(r *http.Request) error {
	if r.Body == nil {
		return errNilBody
	}
	return nil
}

func raise(err error, msg string, code int) exception.Handler {
	return exception.ServerError{
		Error:   err,
		Message: msg,
		Code:    code,
	}
}
