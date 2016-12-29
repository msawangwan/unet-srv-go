package handler

import (
	"net/http"

	"github.com/msawangwan/unet/service/exception"
)

type RouteHandler func(http.ResponseWriter, *http.Request) *exception.Handler

func (rh RouteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := rh(w, r); e != nil {
		http.Error(w, e.Message, e.Code)
	}
}
