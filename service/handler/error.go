package handler

import (
	"errors"

	"net/http"

	"github.com/msawangwan/unet/service/exception"
)

var (
	errNilBody = errors.New("handler: the request body was nil")
)

func checkBody(r *http.Request) error {
	if r.Body == nil {
		return errNilBody
	} else {
		return nil
	}
}

func throw(err error, msg string, code int) *exception.Handler {
	return &exception.Handler{err, msg, code}
}
