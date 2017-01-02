package game

import (
	"time"

	"github.com/msawangwan/unet/env"
)

type Instance struct {
	Key             string `json:"key"`
	ActiveSessionID int    `json:"activeSessionID"`
	Tick            int64  `json:"tick"`
}

func NewInstance(e *env.Global, id string) (*Instance, error) {
	return &Instance{
		Key:             e.CreateKey_GameInstance(id),
		ActiveSessionID: e.AddSession(),
		Tick:            0,
	}, nil
}

func (i *Instance) Start(e *env.Global) error {
	conn, err := e.Get()
	if err != nil {
		return err
	}

	defer func() {
		e.Put(conn)
		e.SetPrefix_Debug()
	}()

	e.SetPrefix("[SESSION START] ")

	conn.Cmd("HSET", i.Key, 0, i.ActiveSessionID)
	conn.Cmd("HSET", i.Key, 1, i.Tick)

	go func() {
		e.Printf("entering loop")
		for {
			e.Printf("top of loop")
			select {
			case <-e.Sessions[i.ActiveSessionID].Terminate:
				e.Printf("terminating session")
			default:
				e.Printf("incrementing tick in loop")
				time.Sleep(150 * time.Millisecond)
				e.Cmd("HSET", i.Key, 1, i.Tick)
			}
		}
	}()

	return nil
}
