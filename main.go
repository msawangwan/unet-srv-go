package main

import (
	"log"

	"net/http"

	"github.com/msawangwan/unet/db"
	"github.com/msawangwan/unet/debug"
	"github.com/msawangwan/unet/env"
	"github.com/msawangwan/unet/service/gateway"
)

func main() {
	var (
		environment *env.Global
	)

	var (
		redis   *db.RedisHandle
		postgre *db.PostgreHandle
		logger  *debug.Log
	)

	var (
		err error
	)

	redis, err = db.NewRedisHandle()
	if err != nil {
		log.Printf("error redis")
	}

	postgre, err = db.NewPostgreHandle()
	if err != nil {
		log.Printf("error pg")
	}

	logger, err = debug.NewLogger()
	if err != nil {
		log.Fatal("error setting up logger") // TODO: fix
	}

	environment = env.New(
		redis,
		postgre,
		logger,
	)

	http.ListenAndServe(":8080", gateway.NewMultiplexer(environment, nil))
}
