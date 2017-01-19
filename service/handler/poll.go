package handler

import (
	"encoding/json"
	"net/http"

	"github.com/msawangwan/unet-srv-go/env"
	"github.com/msawangwan/unet-srv-go/service/exception"

	"github.com/msawangwan/unet-srv-go/engine/game"
)

// PollStart : poll/start
func PollStart(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	j, err := parseJSONInt(r.Body)
	if err != nil {
		return raiseServerError(err)
	}

	gamehandlerstr := game.GameHandlerString(*j)

	sim, err := g.Games.Get(gamehandlerstr)
	if err != nil {
		return raiseServerError(err)
	}

	putconsole := func(s string, id int) {
		g.Prefix("handler", "poller", "start")
		g.Printf("game [%d][%s]", id, s)
		g.PrefixReset()
	}

	putconsole("client waiting for game start ...", *j)

	<-sim.Start // long poller blocking

	putconsole("client got start signal ...", *j)

	// TODO: SEND A HASHED KEY TO ALL CLIENTS temp  willl always send the same debug value

	json.NewEncoder(w).Encode(
		struct {
			Value int `json:"value"`
		}{
			Value: 123456,
		},
	)

	return nil
}

// PollUpdate : poll/update
func PollUpdate(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	return nil
}
