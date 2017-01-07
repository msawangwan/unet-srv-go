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

	GameManager *game.Manager
	KeyGen      *session.KeyGenerator

	sync.Mutex
	sync.WaitGroup
}

// NewGlobalHandle returns a new instance of a global context object
func New(maxSessionsPerHost int, param *config.GameParameters, redis *db.RedisHandle, pg *db.PostgreHandle, log *debug.Log) *Global {
	kgen, err := session.NewKeyGenerator(redis.Pool, log)
	if err != nil {
		log.Printf("env setup error with keygen: %s\n", err.Error()) // TODO: handle for reals
	}

	return &Global{
		GameParameters: param,
		RedisHandle:    redis,
		PostgreHandle:  pg,
		Log:            log,

		GameManager: game.NewGameManager(maxSessionsPerHost),
		KeyGen:      kgen,
	}
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
