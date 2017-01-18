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

	// for below continued see game.Simulation and Table etc.. left off there

	// client sends the game key
	// server concats to make the table look up key
	// client waits
	// server checks redis at interval for number of players
	// when 2 players, server writes to all clients down the channel some sort of starting code?

	return nil
}

// PollUpdate : poll/update
func PollUpdate(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	return nil
}
