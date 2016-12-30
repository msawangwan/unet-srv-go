package handler

import (
	"errors"

	"net/http"
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
