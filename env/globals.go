// package env handles environment and global variables
package env

import (
	"sync"

	"github.com/msawangwan/unet-srv-go/config"
	"github.com/msawangwan/unet-srv-go/db"
	"github.com/msawangwan/unet-srv-go/debug"

	"github.com/msawangwan/unet-srv-go/engine/event"
	"github.com/msawangwan/unet-srv-go/engine/game"
	"github.com/msawangwan/unet-srv-go/engine/manager"
)

// Global encapsulates global handlers
type Global struct {
	*config.GameParameters

	*db.RedisHandle
	*db.PostgreHandle

	*debug.Log

	Games        *game.Table
	KeyManager   *manager.KeyGenerator
	Content      *manager.ContentHandler
	EventEmitter *event.Emitter

	GlobalError chan error

	sync.Mutex
	sync.WaitGroup
}

// New returns a new instance of a global context object
func New(maxSessionsPerHost int, errc chan error, param *config.GameParameters, redis *db.RedisHandle, pg *db.PostgreHandle, log *debug.Log) *Global {
	var (
		lorePath = "lore.json" // TODO: run from config
	)

	checkErr := func(err error) {
		if err != nil {
			defer log.PrefixReset()
			log.Prefix("init", "error")
			log.Fatalf("%s\n", err)
		}
	}

	games, err := game.NewTable(redis.Pool, log)
	checkErr(err)

	kgen, err := manager.NewKeyGenerator(redis.Pool, log)
	checkErr(err)

	cm, err := manager.NewContentHandler(&lorePath)
	checkErr(err)

	emitter, err := event.NewEmitter(redis.Pool, log)
	checkErr(err)

	g := &Global{
		GameParameters: param,

		RedisHandle:   redis,
		PostgreHandle: pg,

		Log: log,

		Games:        games,
		KeyManager:   kgen,
		Content:      cm,
		EventEmitter: emitter,

		GlobalError: errc,
	}

	return g
}
