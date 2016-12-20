package service

import (
	"net/http"
)

type serviceContext struct {
	/* globals here... */
	*Console
}

type serviceHandler struct {
	*serviceContext
	Handle func(*serviceContext, http.ResponseWriter, *http.Request) (int, error)
}

func (sh serviceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if status, err := sh.Handle(sh.serviceContext, w, r); err != nil {
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
