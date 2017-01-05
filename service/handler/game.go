package handler

import (
	//	"time"

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

	dbconn, err := e.Get() // TODO: where to do this? remember to cleanup
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	loop, err = game.CreateNew(skey.RedisFormat, e.GameManager, dbconn, e.Log)
	if err != nil {
		e.Printf("server failed to create new game %s\n", skey.RedisFormat)
	} else {
		e.Printf("server started update routine: %s\n", loop.Label)
	}

	json.NewEncoder(w).Encode(&game.Frame{})

	return nil
}

// POST game/update/frame
func GameFrameUpdate(e *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	return nil
}
