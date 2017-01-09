package main

import (
	"log"
	"os"

	"net/http"

	"github.com/msawangwan/unet-srv-go/config"
	"github.com/msawangwan/unet-srv-go/db"
	"github.com/msawangwan/unet-srv-go/debug"
	"github.com/msawangwan/unet-srv-go/engine/game"
	"github.com/msawangwan/unet-srv-go/env"
	"github.com/msawangwan/unet-srv-go/service/gateway"
)

const (
	kConfigPathFallback = "conf.json" // TODO: use flags instead of args and const
)

var (
	configPath string
	err        error
)

func init() {
	//	runtime.GOMAXPROCS(2)
}

func main() {
	args := os.Args[1:]

	if len(args) == 1 {
		configPath = args[0]
	} else {
		configPath = kConfigPathFallback
	}

	var (
		conf  *config.Configuration
		param *config.GameParameters
	)

	if conf, err = config.LoadConfigurationFile(configPath); err != nil {
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
		make(chan error),
		param,
		redis,
		postgre,
		logger,
	)

	var (
		mainUpdate *game.UpdateHandler
	)

	mainUpdate = game.NewUpdateHandle(environment.GlobalError, environment.Pool, logger)

	go mainUpdate.Monitor()

	logger.Printf("main update loop running ...\n")
	logger.Printf("game engine ready ...\n")
	logger.Printf("game manager ready ...\n")
	logger.Printf("debug logger ready ...\n")
	logger.Printf("redis handle ready ...\n")
	logger.Printf("postgre handle ready ...\n")

	logger.Printf("init done ...\n")
	logger.Printf("all systems go ...\n")

	logger.Printf("service listening and serving on %s ...\n", conf.ListenAddress)

	logger.SetPrefixDefault()

	logger.Fatal(http.ListenAndServe(conf.ListenAddress, gateway.NewMultiplexer(environment, nil)))
}
