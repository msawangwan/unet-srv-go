package service

import (
	"errors"
	"net/http"
)

var (
	ErrNilBody    error = errors.New("nil body read")
	ErrJsonDecode error = errors.New("error decoding json")
	ErrJsonEncode error = errors.New("error encoding json")
)

func nilBodyErr(w http.ResponseWriter, r *http.Request) error {
	if r.Body == nil {
		http.Error(w, "nil body error", 400)
		return ErrNilBody
	} else {
		return nil
	}
}

func jsonDecodeErr(w http.ResponseWriter, err error) {
	http.Error(w, "error decoding json "+err.Error(), 400)
}

func jsonEncodeErr(w http.ResponseWriter, err error) {
	http.Error(w, "error encoding json "+err.Error(), 400)
}
