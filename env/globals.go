// package env handles environment and global variables
package env

import (
	"sync"

	"github.com/msawangwan/unet/config"
	"github.com/msawangwan/unet/db"
	"github.com/msawangwan/unet/debug"
)

// type Global encapsulates global handlers
type Global struct {
	*config.GameParameters
	*db.RedisHandle
	*db.PostgreHandle
	*debug.Log

	Sessions *SessionTable

	sync.Mutex
	sync.WaitGroup
}

// NewGlobalHandle returns a new instance of a global context object
func New(maxSessionsPerHost int, param *config.GameParameters, redis *db.RedisHandle, pg *db.PostgreHandle, log *debug.Log) *Global {
	return &Global{
		GameParameters: param,
		RedisHandle:    redis,
		PostgreHandle:  pg,
		Log:            log,

		Sessions: NewSessionTable(maxSessionsPerHost),
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
