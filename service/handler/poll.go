package handler

import (
	//	"fmt"

	"encoding/json"
	"net/http"

	"github.com/msawangwan/unet-srv-go/env"
	"github.com/msawangwan/unet-srv-go/service/exception"

	"github.com/msawangwan/unet-srv-go/engine/game"
)

type PlayerReadyNotification struct {
	GameID     int
	PlayerName string
}

// PollStart : poll/start
func PollStart(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	var prreq *PlayerReadyNotification

	err := json.NewDecoder(r.Body).Decode(&prreq)
	if err != nil {
		return raiseServerError(err)
	}

	sim, err := g.Games.Get(game.GameLookupString((prreq.GameID)))
	if err != nil {
		return raiseServerError(err)
	}

	sim.PlayerJoinedEvent <- game.OnJoin{prreq.PlayerName}

	select {
	case <-sim.NotifyStart:
		g.Prefix("handler", "pollstart")
		g.Printf("client handler responding to long-poll request, got start signal ...")
		g.PrefixReset()
	case err = <-sim.NotifyError:
		return raiseServerError(err)
	}

	opponent := <-sim.GetOpponent(prreq.PlayerName)

	json.NewEncoder(w).Encode(
		struct {
			Key   int    `json:"key"`
			Value string `json:"value"`
		}{
			Key:   12345, // TODO: non-debug will be random hashed int
			Value: opponent,
		},
	)

	return nil
}

// PollUpdate : poll/update
func PollUpdate(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	return nil
}
