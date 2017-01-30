package handler

import (
	"encoding/json"
	"net/http"

	"github.com/msawangwan/unet-srv-go/engine/game"
	"github.com/msawangwan/unet-srv-go/engine/session"
	"github.com/msawangwan/unet-srv-go/env"
	"github.com/msawangwan/unet-srv-go/service/exception"
)

var labeldebug = func(globals *env.Global) { globals.Label(3, "handler", "game") }

type CreateWorldRequest struct {
	GameKey  int    `json:"key"`
	GameName string `json:"value"`
}

// Loadworld : POST game/world/load
func LoadWorld(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	var (
		createreq *CreateWorldRequest
	)

	err := json.NewDecoder(r.Body).Decode(&createreq)
	if err != nil {
		return raiseServerError(err)
	}

	err = game.LoadWorld(
		g.Content,
		createreq.GameKey,
		g.WorldNodeCount,
		g.WorldScale,
		g.NodeRadius,
		g.MaximumAttemptsWhenSpawningNodes,
		g.Pool,
		g.Log,
	)
	if err != nil {
		return raiseServerError(err)
	}

	//	defer g.PrefixReset()
	//	g.Prefix("handler", "game", "loadworld")
	labeldebug(g)
	defer g.ClearLabel()
	g.Printf("loaded game world [gamekey: %d]", createreq.GameKey)

	err = session.PostGameToLobby(createreq.GameKey, createreq.GameName, g.Pool, g.Log)
	if err != nil {
		return raiseServerError(err)
	}

	return nil
}

type JoinRequest struct {
	GameKey    int    `json:"gameKey"`
	PlayerName string `json:"playerName"`
	Host       bool   `json:"host"`
}

type JoinResponse struct {
	WorldParameters  *game.WorldParameters  `json:"worldParameters"`
	PlayerParameters *game.PlayerParameters `json:"playerParameters"`
}

func JoinGameWorld(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	var (
		joinReq *JoinRequest
	)

	err := json.NewDecoder(r.Body).Decode(&joinReq)
	if err != nil {
		return raiseServerError(err)
	}

	worldparams, err := game.GetWorldParameters(joinReq.GameKey, g.Pool, g.Log)
	if err != nil {
		return raiseServerError(err)
	}

	pid, err := g.KeyManager.GenerateNextPlayerID()
	if err != nil {
		return raiseServerError(err)
	}

	playerparams, err := game.Join(joinReq.GameKey, *pid, joinReq.PlayerName, g.Pool, g.Log)
	if err != nil {
		return raiseServerError(err)
	}

	json.NewEncoder(w).Encode(
		&JoinResponse{
			WorldParameters:  worldparams,
			PlayerParameters: playerparams,
		},
	)

	return nil
}

type CheckNodeHQRequest struct { // TODO: this is a duplicate
	GameID      int    `json:"gameID"`
	PlayerIndex int    `json:"playerIndex"`
	NodeString  string `json:"nodeString"`
}

// game/world/player/hq/validation
func ValidateHQChoice(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	var cnhq *CheckNodeHQRequest

	err := json.NewDecoder(r.Body).Decode(&cnhq)
	if err != nil {
		return raiseServerError(err)
	}

	sim, err := g.Games.Get(game.GameLookupString((cnhq.GameID)))
	if err != nil {
		return raiseServerError(err)
	}

	var isHQValid bool = false

	querychan := sim.CheckNodeValidHQ(cnhq.PlayerIndex, cnhq.NodeString)
	select {
	case b := <-querychan:
		isHQValid = b
	case err = <-sim.NotifyError:
		return raiseServerError(err)
	}

	json.NewEncoder(w).Encode(
		struct {
			Value bool `json:"value"`
		}{
			Value: isHQValid,
		},
	)

	return nil
}

type nodeRequest struct { // TODO: refactor into one, this is duplicated by another struct (see above)
	GameID        int    `json:"gameID"`
	PlayerIndex   int    `json:"playerIndex"`
	NodekeyString string `json:"nodekeyString"`
}

// game/world/node/data
func GetNodeAndCacheData(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	var (
		nr *nodeRequest
	)

	e := json.NewDecoder(r.Body).Decode(&nr)
	if e != nil {
		return raiseServerError(e)
	}

	sim, e := g.Games.Get(game.GameLookupString((nr.GameID)))
	if e != nil {
		return raiseServerError(e)
	}

	labeldebug(g)
	defer g.ClearLabel()
	g.Printf("fetching requested node for caching ... [%s]", nr.NodekeyString)

	c := sim.FetchNodeData(nr.NodekeyString)

	select {
	case nd := <-c:
		labeldebug(g)
		g.Printf("got data [%s], sending it back to [player index: %d]", nr.NodekeyString, nr.PlayerIndex)
		json.NewEncoder(w).Encode(
			struct {
				game.WorldNodeState      `json:"state"`
				game.WorldNodeProperties `json:"properties"`
			}{
				WorldNodeState:      *nd.WorldNodeState,
				WorldNodeProperties: *nd.WorldNodeProperties,
			},
		)
	case e = <-sim.NotifyError:
		return raiseServerError(e)
	}

	return nil
}

// game/world/player/signal/ready
func PlayerSentReadySignal(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	return nil
}
