package main

import (
	"log"
	"runtime"

	"net/http"

	"github.com/msawangwan/unet/config"
	"github.com/msawangwan/unet/db"
	"github.com/msawangwan/unet/debug"
	"github.com/msawangwan/unet/env"
	"github.com/msawangwan/unet/service/gateway"
)

const (
	configFile = "conf.json" // TODO: use flags
)

func main() {
	runtime.GOMAXPROCS(2)

	var (
		err error
	)

	var (
		conf  *config.Configuration
		param *config.GameParameters
	)

	if conf, err = config.LoadConfigurationFile(configFile); err != nil {
		log.Fatal("error loading configuration file:", err.Error())
	}

	if param, err = config.LoadGameParameters(conf.GameParametersFile); err != nil {
		log.Fatal("error loading game parameters file:", err.Error())
	}

	var (
		logger *debug.Log
	)

	logger, err = debug.NewLogger(conf.LogFile)
	if err != nil {
		log.Fatal("error setting up logger")
	}

	logger.SetPrefix_Init()

	var (
		redis   *db.RedisHandle
		postgre *db.PostgreHandle
	)

	redis, err = db.NewRedisHandle()
	if err != nil {
		logger.Fatal("error setting up redis:", err.Error())
	}

	postgre, err = db.NewPostgreHandle()
	if err != nil {
		logger.Fatal("error setting up postgre:", err.Error())
	}

	var (
		environment *env.Global
	)

	environment = env.New(
		conf.MaxGameSessionsAllowed,
		param,
		redis,
		postgre,
		logger,
	)

	logger.Printf("debug logger ready ...\n")
	logger.Printf("redis handle ready ...\n")
	logger.Printf("postgre handle ready ...\n")
	logger.Printf("all systems go ...\n")
	logger.Printf("service listening and serving on %s ...\n", conf.ListenAddress)

	logger.SetPrefix_Debug()

	logger.Fatal(http.ListenAndServe(conf.ListenAddress, gateway.NewMultiplexer(environment, nil)))
}
