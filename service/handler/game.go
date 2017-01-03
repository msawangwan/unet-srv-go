package handler

import (
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

	var (
		gu *game.Update
	)

	gu, err = game.NewInstance(e, skey.RedisFormat)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	go gu.Start(e)

	json.NewEncoder(w).Encode(&game.Frame{})

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
