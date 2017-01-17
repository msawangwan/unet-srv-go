package handler

import (
	//"io"

	"net/http"

	"github.com/msawangwan/unet-srv-go/env"
	"github.com/msawangwan/unet-srv-go/service/exception"
)

// PollStart : poll/start
func PollStart(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	//io.WriteString(w, <-event)

	return nil
}

// PollUpdate : poll/update
func PollUpdate(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	return nil
}
