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
	)

	err := json.NewDecoder(r.Body).Decode(&skey)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	//gl := game.UpdateLoop{skey.RedisFormat}
	//game.UpdateQueue <- gl
	//go game.GameLoop(skey.RedisFormat)
	g := game.NewUpdateInstance(skey.RedisFormat, e.Log)
	go g.OnTick()
	go func() {
		time.Sleep(10 * time.Second)
		g.OnDestroy()
	}()
	e.Printf("game loop queued: %s\n", skey.RedisFormat)

	//var (
	//	gu *game.Update
	//)//

	//gu, err = game.NewInstance(e, skey.RedisFormat)
	//if err != nil {
	//	return &exception.Handler{err, err.Error(), 500}
	//}

	json.NewEncoder(w).Encode(&game.Frame{})
	//r.Close = true
	//go gu.Start(e)

	return nil
}

// POST game/update/frame
func GameFrameUpdate(e *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {

	return nil
}

//func EndGameFrameUpdate(e *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
//	var (
//		skey *session.Key
//	)
//
//	err := json.NewDecoder(r.Body).Decode(&skey)
//	if err != nil {
//		return &exception.Handler{err, err.Error(), 500}
//	}

//	e.SessionTable
