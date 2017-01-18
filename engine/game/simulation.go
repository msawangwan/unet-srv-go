package game

import (
	"errors"
	"sync"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet-srv-go/debug"
)

type Simulation struct {
	Start chan int
	End   chan int
	Turn  chan int
}

func (s *Simulation) StartOnClientSignal(p *pool.Pool, l *debug.Log) error {

	return nil
}

func NewSimulation() (*Simulation, error) {
	return &Simulation{
		Start: make(chan int),
		End:   make(chan int),
		Turn:  make(chan int),
	}, nil
}

type Table struct {
	active map[string]*Simulation
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

func (t *Table) Add(key string) error {
	t.Lock()
	defer func() {
		t.Unlock()
		t.PrefixReset()
	}()

	t.Prefix("game", "simulation")
	t.Printf("adding new simulation [gamename: (key) %s]", key)

	_, exists := t.active[key]
	if exists {
		return errors.New("already added")
	}

	sim, err := NewSimulation()
	if err != nil {
		return err
	}

	t.active[key] = sim

	return nil
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
