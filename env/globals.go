// package env handles environment and global variables
package env

import (
	"sync"

	"github.com/msawangwan/unet/config"
	"github.com/msawangwan/unet/db"
	"github.com/msawangwan/unet/debug"

	"github.com/msawangwan/unet/engine/game"
	"github.com/msawangwan/unet/engine/session"
)

// type Global encapsulates global handlers
type Global struct {
	*config.GameParameters
	*db.RedisHandle
	*db.PostgreHandle
	*debug.Log

	GameManager          *game.Manager
	SessionHandleManager *session.HandleManager
	KeyGen               *session.KeyGenerator

	sync.Mutex
	sync.WaitGroup
}

// NewGlobalHandle returns a new instance of a global context object
func New(maxSessionsPerHost int, param *config.GameParameters, redis *db.RedisHandle, pg *db.PostgreHandle, log *debug.Log) *Global {
	hmanager, err := session.NewHandleManager(100, redis.Pool, log) // TODO: get max capactiy from conf
	checkErr(err, log)

	kgen, err := session.NewKeyGenerator(redis.Pool, log)
	checkErr(err, log)

	g := &Global{
		GameParameters: param,
		RedisHandle:    redis,
		PostgreHandle:  pg,
		Log:            log,

		GameManager:          game.NewGameManager(maxSessionsPerHost),
		SessionHandleManager: hmanager,
		KeyGen:               kgen,
	}

	return g
}

// NullGlobalHanldle  returns an empty global context
func Null() *Global {
	return &Global{
		GameParameters: nil,
		RedisHandle:    nil,
		PostgreHandle:  nil,
		Log:            nil,
	}
}

// checkErr is an un-exported helper function for error checking
func checkErr(err error, log *debug.Log) {
	if err != nil {
		defer log.SetPrefixDefault()
		log.SetPrefix("[INIT][ERROR] ")
		log.Fatalf("%s\n", err)
	}
}
