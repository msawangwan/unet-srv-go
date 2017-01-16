package handler

import (
	"errors"
	//"strconv"

	"encoding/json"
	"net/http"

	"github.com/msawangwan/unet-srv-go/engine/game"
	"github.com/msawangwan/unet-srv-go/env"
	"github.com/msawangwan/unet-srv-go/service/exception"
)

// Loadworld : POST game/world/load
func LoadWorld(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	var (
		gameKey int
	)

	j, err := parseJSONInt(r.Body)
	if err != nil {
		return raiseServerError(err)
	} else if j == nil {
		return raiseServerError(errors.New("nil game key was sent by the client"))
	}

	gameKey = *j

	err = game.LoadWorld(
		gameKey,
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
	g.Printf("loaded game world [gamekey: %d]", gameKey)

	return nil
}

type jbool struct {
	Key   int  `json:"key"`
	Value bool `json:"value"`
}

func JoinGameWorld(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	var (
		j jbool
	)

	err := json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		return raiseServerError(err)
	}

	seed, err := game.GetSeed(j.Key, g.Pool, g.Log)
	if err != nil {
		return raiseServerError(err)
	}

	json.NewEncoder(w).Encode(
		struct {
			Value int64 `json:"value"`
		}{
			Value: *seed,
		},
	)

	return nil
}
