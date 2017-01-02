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

	Sessions           []*ActiveSession
	sessionCount       int
	maxSessionsAllowed int

	sync.Mutex
}

type ActiveSession struct {
	ID        int
	Terminate chan bool
}

// NewGlobalHandle returns a new instance of a global context object
func New(maxSessions int, param *config.GameParameters, redis *db.RedisHandle, pg *db.PostgreHandle, log *debug.Log) *Global {
	return &Global{
		GameParameters:     param,
		RedisHandle:        redis,
		PostgreHandle:      pg,
		Log:                log,
		Sessions:           make([]*ActiveSession, maxSessions),
		sessionCount:       -1,
		maxSessionsAllowed: maxSessions,
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

func (g *Global) AddSession() int {
	g.Lock()
	defer g.Unlock()

	if g.sessionCount > g.maxSessionsAllowed {
		return -1 // TODO: handle
	}

	g.sessionCount += 1

	g.Sessions[g.sessionCount] = &ActiveSession{
		ID:        g.sessionCount,
		Terminate: make(chan bool),
	}

	return g.sessionCount
}

func (g *Global) RemoveSession(id int) {
	g.Lock()
	defer g.Unlock()

	if id < 0 || id > g.maxSessionsAllowed {
		return
	}

	g.Sessions[id] = nil // TODO: kill through channels
	g.sessionCount -= 1

	return
}
