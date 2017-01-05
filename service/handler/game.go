package handler

import (
	"time"

	"encoding/json"
	"net/http"

	"github.com/msawangwan/unet/engine/game"
	"github.com/msawangwan/unet/engine/session"
	"github.com/msawangwan/unet/env"
	"github.com/msawangwan/unet/service/exception"
)

// POST game/update/start
func StartGameUpdate(e *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	var (
		skey *session.Key
		loop *game.Update
	)

	err := json.NewDecoder(r.Body).Decode(&skey)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	loop = game.NewUpdateInstance(skey.RedisFormat, e.Log)

	go loop.OnTick()
	go func() { // TODO: debug only
		time.Sleep(30 * time.Second)
		loop.OnDestroy()
	}()

	e.Printf("server started a new game session: %s\n", skey.RedisFormat)

	json.NewEncoder(w).Encode(&game.Frame{})

	return nil
}

// POST game/update/frame
func GameFrameUpdate(e *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	return nil
}
