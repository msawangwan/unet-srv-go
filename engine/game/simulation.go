package game

import (
	"errors"
	"sync"
	"time"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet-srv-go/debug"
)

type Simulation struct {
	handlerstring string

	Start chan bool
	End   chan bool
	Turn  chan int
	Error chan error

	StartTick  *time.Ticker
	UpdateTick *time.Ticker
}

func (s *Simulation) WaitUntilAllClientsReady(p *pool.Pool, l *debug.Log) chan bool {
	consolePrefix := func() { l.Prefix("game", "simulation") }

	// commented out cause the compiler complains if it's not used.. consider
	// converting into its own function?
	//	onStart := func() {
	//		s.Start <- true
	//		close(s.Start)
	//	}

	go func() {
		for {
			select {
			case <-s.StartTick.C:
				consolePrefix()
				l.Printf("waiting for players [game: %s][start ticker]", s.handlerstring)
				// check redis for player count ...
				// if all players ready then call onStart()
				l.PrefixReset()
			case <-s.Start:
				consolePrefix()
				l.Printf("all players joined, starting game")
				l.PrefixReset()

				s.StartTick.Stop()

				return
			}
		}
	}()

	return make(chan bool) // future clients will read fom this channel to get a ready signal TODO: add this to a table??
}

func NewSimulation(handlerstring string) (*Simulation, error) {
	return &Simulation{
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
	// add another map that returns the start signal channel??
	// add another map that returns the update/turn signal channel??
	sync.Mutex
	*pool.Pool
	*debug.Log
}

func NewTable(p *pool.Pool, l *debug.Log) (*Table, error) {
	return &Table{
		active: make(map[string]*Simulation),
		Pool:   p,
		Log:    l,
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

	sim, err := NewSimulation(key)
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
