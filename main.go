package main

import (
	"log"

	"net/http"

	"github.com/msawangwan/unet/db"
	"github.com/msawangwan/unet/debug"
	"github.com/msawangwan/unet/env"
	"github.com/msawangwan/unet/service/gateway"
)

var (
	laddr = ":8080" // TODO: get from conf
)

func main() {
	var (
		err error
	)

	var (
		logger *debug.Log
	)

	logger, err = debug.NewLogger()
	if err != nil {
		log.Fatal("error setting up logger")
	} else {
		logger.Printf("debug logger ready ...\n")
	}

	var (
		redis   *db.RedisHandle
		postgre *db.PostgreHandle
	)

	redis, err = db.NewRedisHandle()
	if err != nil {
		logger.Fatal("error setting up redis")
	} else {
		logger.Printf("redis handle ready ...\n")
	}

	postgre, err = db.NewPostgreHandle()
	if err != nil {
		logger.Fatal("error pg")
	} else {
		logger.Printf("postgre handle ready ...\n")
	}

	var (
		environment *env.Global
	)

	environment = env.New(
		redis,
		postgre,
		logger,
	)

	logger.Printf("all systems go ...\n")
	logger.Printf("service listening and serving on %s ...\n", laddr)
	logger.Fatal(http.ListenAndServe(laddr, gateway.NewMultiplexer(environment, nil)))
}
