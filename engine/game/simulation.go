package game

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"

	"github.com/msawangwan/unet-srv-go/debug"
	"github.com/msawangwan/unet-srv-go/engine/prng"
)

type Simulation struct {
	gametable *Table

	lookupstr string

	//PlayerLeft        chan struct{}
	PlayerJoinedEvent      chan OnJoin
	BroadcastPlayerToAct   chan OnTurn
	SendTurnCompletedEvent chan OnTurn

	NotifyStart chan struct{}
	NotifyEnd   chan struct{}
	NotifyTurn  chan struct{}
	NotifyError chan error

	SendMessage map[string]chan string // TODO: come up with better structure

	StartTick  *time.Ticker
	UpdateTick *time.Ticker

	playercount int
	gameslots   [10]string

	putconsole func(stdout string)

	rand *prng.Manager

	*pool.Pool
	*debug.Log
}

func NewSimulation(gametable *Table, lookupstr string, p *pool.Pool, l *debug.Log) (*Simulation, error) {
	const (
		maxplayerspergame = 2 // TODO: load from config.json or params.json
	)

	l.Label(1, "game", "newsimulation")
	defer l.PrefixReset()

	l.Printf("instantiating a new simulation ...")

	worldseed, err := GetWorldSeed(lookupstr, p, l)
	if err != nil {
		return nil, err
	}

	rand, err := prng.NewInstanceManager(maxplayerspergame, *worldseed)
	if err != nil {
		return nil, err
	}

	sim := &Simulation{
		gametable: gametable,

		lookupstr: lookupstr,

		PlayerJoinedEvent:      make(chan OnJoin),
		BroadcastPlayerToAct:   make(chan OnTurn),
		SendTurnCompletedEvent: make(chan OnTurn),

		NotifyStart: make(chan struct{}),
		NotifyEnd:   make(chan struct{}),
		NotifyTurn:  make(chan struct{}),
		NotifyError: make(chan error),

		StartTick:  time.NewTicker(750 * time.Millisecond),
		UpdateTick: time.NewTicker(2500 * time.Millisecond),

		SendMessage: make(map[string]chan string),

		playercount: 0,

		rand: rand,

		Pool: p,
		Log:  l,
	}

	putconsole := func(msg string) {
		sim.Prefix("simulation", sim.lookupstr)
		sim.Printf("[%s] %s", sim.lookupstr, msg)
		sim.PrefixReset()
	}

	sim.putconsole = putconsole

	sim.SendCurrentPlayerTurn()

	return sim, nil
}

func (s *Simulation) WaitUntilAllClientsReady() (chan struct{}, error) { // TODO: rename to ALLplayerjoined
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

	return s.NotifyStart, nil
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

// SendCurrentPlayerTurn will, every interval of x, send down a channel the
// current player's turn, anyone subscribed will get this message (note:
// potential bugs here as once the message is read it is discarded)
func (s *Simulation) SendCurrentPlayerTurn() {
	go func() {
		conn, err := s.Get()
		if err != nil {
			s.NotifyError <- err
			return
		}
		defer s.Put(conn)

		current, err := conn.Cmd("HGET", s.lookupstr, "game_player_to_act").Int()
		if err != nil {
			s.NotifyError <- err
			return
		}

		for {
			select {
			case <-s.UpdateTick.C:
				s.putconsole("update tick...")
				s.putconsole(fmt.Sprintf("player to act [%d]", current))
				current, err = conn.Cmd("HGET", s.lookupstr, "game_player_to_act").Int()
				if err != nil {
					s.NotifyError <- err // TODO: break or continue??
				}

				s.BroadcastPlayerToAct <- OnTurn{PlayerToAct: current}
			}
		}
	}()
}

func (s *Simulation) NotifyTurnComplete(playerindex int) {
	go func() {
		var next int = -1
		if playerindex == 1 {
			next = 2
		} else {
			next = 1
		}

		if next == -1 {
			s.NotifyError <- errors.New("wtf cant have -1 player index")
		}

		s.SendTurnCompletedEvent <- OnTurn{PlayerToAct: next} // TODO: create a new type for this
		s.Cmd("HSET", s.lookupstr, "game_player_to_act", next)

		s.putconsole(fmt.Sprintf("player [%d] sent turn complete to server, it is now player [%d] turn", playerindex, next))

		return
	}()
}

func (s *Simulation) NotifyPlayerTurnStart(playerindex int) chan struct{} {
	onStart := make(chan struct{})

	go func() {
		for {
			select {
			case curr := <-s.BroadcastPlayerToAct:
				s.putconsole("read broadcast")
				if curr.PlayerToAct == playerindex {
					onStart <- struct{}{}
					return
				} else {
					s.putconsole("not me so continue")
				}
			case <-s.SendTurnCompletedEvent:
				s.putconsole("got a turn completed event so bailing out")
				onStart <- struct{}{}
				return
			}
		}
	}()

	return onStart
}

func (s *Simulation) CheckNodeValidHQ(playerindex int, nodestr string) chan bool {
	sendvalid := make(chan bool)

	s.putconsole("verify node is valid hq")

	go func() { // check for truncated 0's after the decimal! we do this on the client only right now
		nodestatstr := GameNodeStatString(s.lookupstr, nodestr)
		s.putconsole(fmt.Sprintf("checking if node is hq [%s][%s]", s.lookupstr, nodestatstr))
		b := s.Cmd("HGET", nodestatstr, "node_ishq")
		if b.Err != nil {
			s.NotifyError <- b.Err
			return
		} else if b.IsType(redis.Nil) {
			s.putconsole("that's not a valid node")
			sendvalid <- false
			return
		}

		result, err := b.Str()
		if err != nil {
			s.NotifyError <- err
		}

		s.putconsole(fmt.Sprintf("result from server (before parsing) [%s]", result))
		isHQ, err := strconv.ParseBool(result)
		if err != nil {
			s.NotifyError <- err
			return
		}
		if isHQ {
			s.putconsole("node is aleady hq")
			s.putconsole("player cannot choose this node as their hq")
			sendvalid <- false
		} else {
			s.putconsole("node has not been selected as an hq")
			s.putconsole("assigning node as hq to player")
			conn, err := s.Get()
			if err != nil {
				s.NotifyError <- err
			}
			defer s.Put(conn)

			if err = conn.Cmd("MULTI").Err; err != nil {
				s.NotifyError <- err
				return
			}

			conn.Cmd("HSET", nodestatstr, "node_ishq", strconv.FormatBool(true))
			conn.Cmd("HSET", nodestatstr, "node_hq_owner_by_index", playerindex)
			conn.Cmd("HSET", GamePlayerNString(s.lookupstr, playerindex), "player_hq", nodestatstr)

			if err = conn.Cmd("EXEC").Err; err != nil {
				s.NotifyError <- err
				return
			}

			sendvalid <- true
		}
		close(sendvalid)
		return
	}()

	return sendvalid
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
