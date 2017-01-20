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

	lookupstr string

	//PlayerLeft        chan struct{}
	PlayerJoinedEvent chan OnJoin

	NotifyStart chan struct{}
	NotifyEnd   chan struct{}
	NotifyTurn  chan int
	NotifyError chan error

	SendMessage map[string]chan string // TODO: come up with better structure

	StartTick  *time.Ticker
	UpdateTick *time.Ticker

	playercount int
	gameslots   [10]string

	putconsole func(stdout string)

	*pool.Pool
	*debug.Log
}

func (s *Simulation) WaitUntilAllClientsReady() (chan struct{}, error) {
	go func() {
		s.putconsole("new game started")

		for {
			select {
			case <-s.StartTick.C:
				s.putconsole("waiting for opponent to join ...")

				if s.playercount > 1 {
					go func() { // kill routine
						s.StartTick.Stop()
						s.NotifyStart <- struct{}{}
						close(s.NotifyStart)
					}()
				}
			case event := <-s.PlayerJoinedEvent:
				s.putconsole(fmt.Sprintf("player joined [%s]", event.PlayerName))
				s.playercount += 1
				s.gameslots[s.playercount-1] = event.PlayerName
			case <-s.NotifyStart:
				s.putconsole("all players joined, notifying opponent and starting game")
				s.putconsole("printing player list")

				playerlist, err := s.Cmd("SMEMBERS", GamePlayerListString(s.lookupstr)).List()
				if err != nil {
					s.NotifyError <- err
				}

				for i, player := range playerlist {
					s.putconsole(fmt.Sprintf("%d) %s", i, player))
				}

				return
			}
		}
	}()

	return s.NotifyStart, nil // future clients will read fom this channel to get a ready signal??
}

func (s *Simulation) GetOpponent(playername string) chan string {
	sendmsg := make(chan string)

	s.putconsole("finding opponent name ...")

	go func() {
		for _, name := range s.gameslots {
			if name == playername {
				continue
			}
			s.putconsole(fmt.Sprintf("found oppoent [%s]", playername))
			sendmsg <- name
			return
		}
	}()

	return sendmsg
}

func NewSimulation(gametable *Table, lookupstr string, p *pool.Pool, l *debug.Log) (*Simulation, error) {
	sim := &Simulation{
		gametable: gametable,

		lookupstr: lookupstr,

		PlayerJoinedEvent: make(chan OnJoin),

		NotifyStart: make(chan struct{}),
		NotifyEnd:   make(chan struct{}),
		NotifyTurn:  make(chan int),
		NotifyError: make(chan error),

		StartTick:  time.NewTicker(750 * time.Millisecond),
		UpdateTick: time.NewTicker(2500 * time.Millisecond),

		SendMessage: make(map[string]chan string),

		playercount: 0,

		Pool: p,
		Log:  l,
	}

	putconsole := func(msg string) {
		sim.Prefix("simulation", sim.lookupstr)
		sim.Printf("[%s] %s", sim.lookupstr, msg)
		sim.PrefixReset()
	}

	sim.putconsole = putconsole

	return sim, nil
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
	t.Printf("adding new simulation [gamelookupstring: %s]", key)

	_, exists := t.active[key]
	if exists {
		return nil, errors.New("already added")
	}

	sim, err := NewSimulation(t, key, t.Pool, t.Log)
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
	t.Printf("accessing simulation [gamelookupstring: %s]", key)

	_, exists := t.active[key]
	if exists {
		t.Printf("found simulation")
		return t.active[key], nil
	}

	t.Printf("no simulation with that key found")
	return nil, errors.New(fmt.Sprintf("no simulation has lookup string [%s]", key))
}
