// package env handles environment and global variables
package env

import (
	"github.com/msawangwan/unet/db"
	"github.com/msawangwan/unet/debug"
)

// type Global encapsulates global handlers
type Global struct {
	*db.RedisHandle
	*db.PostgreHandle
	*debug.Log
}

// NewGlobalHandle returns a new instance of a global context object
func New(redis *db.RedisHandle, pg *db.PostgreHandle, log *debug.Log) *Global {
	return &Global{
		RedisHandle:   redis,
		PostgreHandle: pg,
		Log:           log,
	}
}

// NullGlobalHanldle  returns an empty global context
func Null() *Global {
	return &Global{
		RedisHandle:   nil,
		PostgreHandle: nil,
		Log:           nil,
	}
}
