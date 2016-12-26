// package env handles environment and global variables
package env

import (
	"github.com/msawangwan/unet/db"
)

// type Global encapsulates global handlers
type Global struct {
	*db.RedisHandle
	*db.PostgreHandle
}

func NewGlobalHandle(redis *db.RedisHandle, pg *db.PostgreHandle) *Global {
	return &Global{
		RedisHandle:   redis,
		PostgreHandle: pg,
	}
}

func NullGlobalHandle() *Global {
	return &Global{
		RedisHandle:   nil,
		PostgreHandle: nil,
	}
}
