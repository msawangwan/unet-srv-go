package handler

import (
	"errors"

	"encoding/json"
	"net/http"

	//"github.com/msawangwan/unet-srv-go/engine/session"
	"github.com/msawangwan/unet-srv-go/env"
	"github.com/msawangwan/unet-srv-go/service/exception"
)

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

	defer g.PrefixReset()
	g.Prefix("handler", "game", "loadworld")
	g.Printf("loading game world [gamekey: %d]", gameKey)

	json.NewEncoder(w).Encode(
		struct{}{},
	)

	return nil
}
