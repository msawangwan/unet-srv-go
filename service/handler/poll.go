package handler

import (
	"encoding/json"
	"net/http"

	"github.com/msawangwan/unet-srv-go/env"
	"github.com/msawangwan/unet-srv-go/service/exception"

	"github.com/msawangwan/unet-srv-go/engine/game"
)

type PlayerReadyNotification struct {
	GameID     int    `json:"gameID"`
	PlayerName string `json:"playerName"`
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

type PlayerTurnPollRequest struct {
	GameID      int `json:"gameID"`
	PlayerIndex int `json"playerIndex"`
}

// game/turn/poll
func PollForTurnSignal(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	var ptpr *PlayerTurnPollRequest

	err := json.NewDecoder(r.Body).Decode(&ptpr)
	if err != nil {
		return raiseServerError(err)
	}

	gamelookupstr := game.GameLookupString(ptpr.GameID)

	sim, err := g.Games.Get(gamelookupstr)
	if err != nil {
		return raiseServerError(err)
	}

	<-sim.NotifyPlayerTurnStart(ptpr.PlayerIndex) // block until we get note to go

	return nil
}

type PlayerNotifyServer struct {
	GameID      int `json:"gameID"`
	PlayerIndex int `json:"playerIndex"`
}

// game/turn/complete
func GotPlayerTurnComplete(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	var pns *PlayerNotifyServer
	err := json.NewDecoder(r.Body).Decode(&pns)
	if err != nil {
		return raiseServerError(err)
	}

	gamelookupstr := game.GameLookupString(pns.GameID)

	sim, err := g.Games.Get(gamelookupstr)
	if err != nil {
		return raiseServerError(err)
	}

	sim.NotifyTurnComplete(pns.PlayerIndex)

	return nil
}

// PollUpdate : poll/update
func PollUpdate(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	return nil
}
