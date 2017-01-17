package handler

import (
	"encoding/json"
	"net/http"

	"github.com/msawangwan/unet-srv-go/engine/game"
	"github.com/msawangwan/unet-srv-go/engine/session"
	"github.com/msawangwan/unet-srv-go/env"
	"github.com/msawangwan/unet-srv-go/service/exception"
)

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

	defer g.PrefixReset()
	g.Prefix("handler", "game", "loadworld")
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

func JoinGameWorld(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	var (
		joinReq *JoinRequest
	)

	err := json.NewDecoder(r.Body).Decode(&joinReq)
	if err != nil {
		return raiseServerError(err)
	}

	//seed, err := game.GetSeed(joinReq.GameKey, g.Pool, g.Log)
	//if err != nil {
	//	return raiseServerError(err)
	//}

	params, err := game.GetGameParameters(joinReq.GameKey, g.Pool, g.Log)
	if err != nil {
		return raiseServerError(err)
	}

	game.Join(joinReq.GameKey, joinReq.PlayerName, g.Pool, g.Log)

	//json.NewEncoder(w).Encode(
	//	struct {
	//		Value int64 `json:"value"`
	//	}{
	//		Value: *seed,
	//	},
	//)

	json.NewEncoder(w).Encode(params)

	return nil
}
