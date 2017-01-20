package event

import (
	//"fmt"
	"strconv"

	"github.com/msawangwan/unet-srv-go/debug"

	"github.com/mediocregopher/radix.v2/pool"
)

type Event struct {
	Message string
}

type Emitter struct {
	Registrar map[int]chan Event

	*pool.Pool
	*debug.Log
}

func NewEmitter(p *pool.Pool, l *debug.Log) (*Emitter, error) {
	e := &Emitter{
		Registrar: make(map[int]chan Event),
		Pool:      p,
		Log:       l,
	}
	return e, nil
}

func (e *Emitter) Register(gamekey int, playername string) error {
	conn, err := e.Get()
	if err != nil {
		return err
	}

	defer func() {
		e.Put(conn)
		e.PrefixReset()
	}()

	conn.Cmd("HGET", "hk_gameHnadleKey", gamekeyString(gamekey)).Str()

	return nil
}

func gamekeyString(gamekey int) string {
	return strconv.Itoa(gamekey)
}
