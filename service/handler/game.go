package handler

import (
	"errors"

	//"encoding/json"
	"net/http"

	"github.com/msawangwan/unet-srv-go/engine/game"
	"github.com/msawangwan/unet-srv-go/env"
	"github.com/msawangwan/unet-srv-go/service/exception"
)

// Loadworld : POST game/load/world
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
