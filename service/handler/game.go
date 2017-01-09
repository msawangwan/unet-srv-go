package handler

import (
	"encoding/json"
	"net/http"

	"github.com/msawangwan/unet-srv-go/env"

	"github.com/msawangwan/unet-srv-go/engine/game"
	"github.com/msawangwan/unet-srv-go/engine/session"

	"github.com/msawangwan/unet-srv-go/service/exception"
)

// POST game/update/start
func StartGameUpdate(e *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	var (
		skey *session.Key
		loop *game.Update
	)

	err := json.NewDecoder(r.Body).Decode(&skey)
	if err != nil {
		return raiseServerError(err)
	}

	rkey := e.CreateHashKey_SessionGameUpdateLoop(skey.RedisFormat) // TODO: rename to game state or something better

	loop, err = game.CreateNew(skey.RedisFormat, rkey, e.GameManager, e.Pool, e.Log)
	if err != nil {
		e.Printf("server failed to create new game %s\n", skey.RedisFormat)
	} else {
		e.Printf("server started update routine: %s\n", loop.Label)
	}

	json.NewEncoder(w).Encode(&game.Frame{})

	return nil
}

// // JoinGame : POST
// func JoinGame(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {

// 	return nil
// }

// POST game/update/enter
func EnterGameUpdate(e *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	var (
		skey *session.Key
		loop *game.Update
	)

	err := json.NewDecoder(r.Body).Decode(&skey)
	if err != nil {
		return raise(err, "MISSING KEY", 500)
	}

	//	rkey := e.CreateHashKey_SessionGameUpdateLoop(skey.RedisFormat)

	loop, err = game.EnterExisting(skey.RedisFormat, e.GameManager, e.Log)
	if err != nil {
		return raise(err, err.Error(), 500)
	}

	e.Printf("joined an existing game on the server: %s\n", loop.Label)

	json.NewEncoder(w).Encode(&game.Frame{})

	return nil
}

// POST game/update/kill
func KillGameUpdate(e *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	var (
		skey *session.Key
	)

	// TODO: need robust validation

	err := json.NewDecoder(r.Body).Decode(&skey)
	if err != nil {
		raiseServerError(err)
	}

	rkey := e.CreateHashKey_SessionGameUpdateLoop(skey.RedisFormat) // TODO: rename to game state or something better

	_, err = game.EndActive(skey.RedisFormat, rkey, e.GameManager, e.Log)
	if err != nil {
		raiseServerError(err)
	}

	return nil
}

// POST game/update/frame
func GameFrameUpdate(e *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	return nil
}
