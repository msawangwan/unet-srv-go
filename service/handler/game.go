package handler

import (
	"net/http"

	"github.com/msawangwan/unet/engine/game"
	"github.com/msawangwan/unet/env"
	"github.com/msawangwan/unet/service/exception"
)

func StartGame(e *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	var (
		gi *game.Instance
	)

	gi, err := game.NewInstance(e, "")
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}
	gi.Start(e)
	// return info to track the goroutine
	return nil
}

func GameFrameUpdate(e *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {

	return nil
}
