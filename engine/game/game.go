package game

import (
	"sync"
	"time"

	"encoding/json"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet/debug"
)

const (
	kMaxDuration = 1 * time.Hour // more of a saftey, or for special events, increase or remove later
	kTimeout     = 5 * time.Minute
	kTick        = 2500 * time.Millisecond
)

const (
	kMaxPlayers = 2
)

type Update struct {
	Label       string `json:"label"`
	InstanceKey string
	playerCount int

	Players [10]string

	Timer  *time.Timer
	Ticker *time.Ticker

	Error chan error
	Done  chan bool

	*pool.Pool
	*debug.Log

	sync.Mutex
}

func NewUpdateRoutine(label string, key string, conns *pool.Pool, log *debug.Log) (*Update, error) {
	update := &Update{
		Label:       label,
		InstanceKey: key,
		Timer:       time.NewTimer(kMaxDuration),
		Ticker:      time.NewTicker(kTick),
		Error:       make(chan error),
		Done:        make(chan bool),
		Pool:        conns,
		Log:         log,
	}

	conn, err := update.Get()
	if err != nil {
		return nil, err
	}
	defer update.Put(conn)

	if err := conn.Cmd("MULTI").Err; err != nil {
		return nil, err
	}

	conn.Cmd("HSET", update.InstanceKey, "game-name", update.Label)
	conn.Cmd("HSET", update.InstanceKey, "frame", 0)

	if err := conn.Cmd("EXEC").Err; err != nil {
		return nil, err
	}

	return update, nil
}

func (u *Update) Enter(player string) error {
	if u.playerCount >= kMaxPlayers {
		u.Printf("more than 2 players") // TODO: handle
	}

	conn, err := u.Get()
	if err != nil {
		return err
	}
	defer u.Put(conn)

	if err := conn.Cmd("MULTI").Err; err != nil {
		return err
	}

	var all []byte

	u.Lock()
	{
		n := len(u.Players)
		if (n + 1) <= cap(u.Players) {
			u.Players[n] = "a player" + string(n)
		}

		//u.Players = append(u.Players, player)
		all, err = json.Marshal(u.Players)
		u.Printf("players as byte arr: %v\n", all)
		u.Printf("players as slice: %s\n", string(all))
		u.Printf("players as slice: %s\n", string(all[:]))
	}
	u.Unlock()

	conn.Cmd("HSET", u.InstanceKey, "player-count", u.playerCount)
	//	conn.Cmd("HSET", u.InstanceKey, "players", string(all[:]))

	if err := conn.Cmd("EXEC").Err; err != nil {
		return err
	}

	return nil
}

func (u *Update) OnTick() {
	conn, err := u.Get()
	if err != nil {
		u.sendErr(err)
		return
	}
	defer u.Put(conn)

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

			conn.Cmd("HINCRBY", u.InstanceKey, "frame", 1)
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

func (u *Update) sendErr(err error) {
	u.Error <- err
	close(u.Error)
}
