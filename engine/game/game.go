package game

import (
	"time"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet/debug"
)

const (
	kMaxDuration = 1 * time.Hour // more of a saftey, or for special events, increase or remove later
	kTimeout     = 5 * time.Minute
	kTick        = 2500 * time.Millisecond
)

type Update struct {
	Label       string `json:"label"`
	InstanceKey string

	Timer  *time.Timer
	Ticker *time.Ticker

	Error chan error
	Done  chan bool

	*pool.Pool
	*debug.Log
}

func NewUpdateRoutine(label string, key string, conns *pool.Pool, log *debug.Log) *Update {
	return &Update{
		Label:       label,
		InstanceKey: key,
		Timer:       time.NewTimer(kMaxDuration),
		Ticker:      time.NewTicker(kTick),
		Error:       make(chan error),
		Done:        make(chan bool),
		Pool:        conns,
		Log:         log,
	}
}

func (u *Update) OnTick() {
	conn, err := u.Get()
	if err != nil {
		// TODO: send down error chan
		u.Printf("%s\n", err.Error())
	}
	defer u.Put(conn)

	if err := conn.Cmd("MULTI").Err; err != nil {
		// TODO: send down error chan
		u.Printf("%s\n", err.Error())
	}

	conn.Cmd("HSET", u.InstanceKey, 0, u.Label)
	conn.Cmd("HSET", u.InstanceKey, 1, 0)

	if err := conn.Cmd("EXEC").Err; err != nil {
		// TODO: err chan
		u.Printf("%s\n", err.Error())
	}

	for {
		select {
		case <-u.Timer.C:
			u.SetPrefix("[UPDATE][ON_TIMEOUT] ")
			u.Printf("timer expired: %s\n", u.Label)
			u.SetPrefixDefault()
		case <-u.Ticker.C:
			u.SetPrefix("[UPDATE][ON_TICK] ")
			u.Printf("tick: %s\n", u.Label)
			u.SetPrefixDefault()

			conn.Cmd("HINCRBY", u.InstanceKey, 1, 1)
		case <-u.Done:
			u.SetPrefix("[UPDATE][ON_DONE] ")
			u.Printf("loop terminated: %s\n", u.Label)
			u.SetPrefixDefault()

			u.Timer.Stop()
			u.Ticker.Stop()

			return
		}
	}
}

func (u *Update) OnDestroy() {
	u.Done <- true
	close(u.Done)
}
