package game

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet-srv-go/debug"
)

type Simulation struct {
	gametable *Table

	handlerstring string

	Start chan bool
	End   chan bool
	Turn  chan int
	Error chan error

	StartTick  *time.Ticker
	UpdateTick *time.Ticker

	playercount int
}

func (s *Simulation) WaitUntilAllClientsReady(p *pool.Pool, l *debug.Log) (chan bool, error) {
	outconsole := func(s string) {
		l.Prefix("game", "simulation")
		l.Printf(s)
	}

	outconsole(fmt.Sprintf("new game started [%s], waiting for joining players", s.handlerstring))

	gameplayerliststr := fmt.Sprintf("%s:%s", s.handlerstring, "playerlist")

	go func() {
		conn, err := p.Get()
		if err != nil {
			s.Error <- err
			return
		}
		defer p.Put(conn)

		for {
			select {
			case <-s.StartTick.C:
				outconsole(fmt.Sprintf("waiting for players [game: %s][start ticker]", s.handlerstring))

				playerlist, err := conn.Cmd("SMEMBERS", gameplayerliststr).List()
				if err != nil {
					s.Error <- err // TODO: return or not?
				}

				if len(playerlist) >= 2 {
					go func() { // kill routine
						outconsole(fmt.Sprintf("player count is [%d], notifying clients", len(playerlist)))
						s.Start <- true
						close(s.Start)
					}()
				}

				l.PrefixReset()
			case <-s.Start:
				outconsole(fmt.Sprintf("all players joined, starting game [%s]", s.handlerstring))

				s.StartTick.Stop()

				l.PrefixReset()
				return
			}
		}
	}()

	return s.Start, nil // future clients will read fom this channel to get a ready signal??
}

func NewSimulation(gametable *Table, handlerstring string) (*Simulation, error) {
	return &Simulation{
		gametable: gametable,

		handlerstring: handlerstring,

		Start: make(chan bool),
		End:   make(chan bool),
		Turn:  make(chan int),
		Error: make(chan error),

		StartTick:  time.NewTicker(750 * time.Millisecond),
		UpdateTick: time.NewTicker(2500 * time.Millisecond),
	}, nil
}

type Table struct {
	active map[string]*Simulation

	waiting map[string]chan bool // start signals
	update  map[string]chan bool // update signals
	turn    map[string]chan bool // might not need this, can just use update or this

	sync.Mutex
	*pool.Pool
	*debug.Log
}

func NewTable(p *pool.Pool, l *debug.Log) (*Table, error) {
	return &Table{
		active: make(map[string]*Simulation),

		waiting: make(map[string]chan bool),
		update:  make(map[string]chan bool),
		turn:    make(map[string]chan bool),

		Pool: p,
		Log:  l,
	}, nil
}

func (t *Table) Add(key string) (*Simulation, error) {
	t.Lock()
	defer func() {
		t.Unlock()
		t.PrefixReset()
	}()

	t.Prefix("game", "simulation")
	t.Printf("adding new simulation [gamename: (key) %s]", key)

	_, exists := t.active[key]
	if exists {
		return nil, errors.New("already added")
	}

	sim, err := NewSimulation(t, key)
	if err != nil {
		return nil, err
	}

	t.active[key] = sim

	return sim, nil
}

func (t *Table) Get(key string) (*Simulation, error) {
	t.Lock()
	defer func() {
		t.Unlock()
		t.PrefixReset()
	}()

	t.Prefix("game", "simulation")
	t.Printf("accessing simulation [gamename (key): %s]", key)

	_, exists := t.active[key]
	if exists {
		t.Printf("found simulation")
		return t.active[key], nil
	}

	t.Printf("no simulation with that key found")
	return nil, errors.New("no simulation with that name")
}
