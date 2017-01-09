package game

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"encoding/json"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet-srv-go/debug"
)

const (
	maxPlayersAllowed = 10
)

const (
	keySimulationInstance = "sim:key"
)

const (
	keySimulationLabel = "label"
	keyFrame           = "frame"
	keyPlayerCount     = "player_count"
	keyPlayerList      = "player_list"
)

const (
	durTimeoutDebug   = 1 * time.Hour
	durTimeout        = 5 * time.Minute
	durUpdateInterval = 2500 * time.Millisecond
)

var (
	errAlreadyExists = errors.New("simulation already exists, failed to create")
)

// Simulation is a game update loop abstraction
type Simulation struct {
	Label string `json:"label"`
	Seed  int64  `json:"seed"`

	Players []string `json:"players"`

	OnComplete chan bool  `json:"-"`
	OnError    chan error `json:"-"`

	Timer  *time.Timer  `json:"-"`
	Ticker *time.Ticker `json:"-"`

	*pool.Pool `json:"-"`
	*debug.Log `json:"-"`

	sync.Mutex `json:"-"`
}

// NewSimulation returns a new simulation instance
func NewSimulation(label string, seed int64, errc error chan, conns *pool.Pool, log *debug.Log) (*Simulation, error) {
	s := &Simulation{
		Label: label,
		Seed:  seed,

		Players: make([]string, 0, maxPlayersAllowed),

		OnComplete: make(chan bool),
        OnError: errc,
		// OnError:    make(chan error),

		Timer:  time.NewTimer(durTimeoutDebug),
		Ticker: time.NewTicker(durUpdateInterval),

		Pool: conns,
		Log:  log,
	}

	conn, err := s.Get()
	if err != nil {
		return nil, err
	}

	defer func() {
		s.Put(conn)
		s.SetPrefixDefault()
	}()

	s.SetPrefix("[GAME][SIMULATION][NEW] ")

	key := makeKey(keySimulationInstance, s.Label)

	check, err := conn.Cmd("EXISTS", key).Int()
	if err != nil {
		return nil, err
	} else if check == 1 {
		return nil, errAlreadyExists
	} else {
		s.Printf("created new simulation: [key: %s] [label: %s] ...", key, s.Label)
		conn.Cmd("HMSET", key, keySimulationLabel, s.Label, keyFrame, 0)
	}

	return s, nil
}

// Join is called by the client to register to simulation updates
func (s *Simulation) Join(player string) error {
	conn, err := s.Get()
	if err != nil {
		return err
	}

	s.Lock()

	defer func() {
		s.Unlock()
		s.Put(conn)
		s.SetPrefixDefault()
	}()

	s.SetPrefix("[GAME][SIMULATION][JOIN] ")

	key := makeKey(keySimulationInstance, s.Label)

	n := len(s.Players)
	if n >= maxPlayersAllowed {
		s.Printf("more than 2 players") // TODO: handle
	}

	var all []byte

	if (n + 1) <= cap(s.Players) {
		s.Players = append(s.Players, player)
	} else {
		s.Printf("failed to add player: %s", player)
	}

	all, err = json.Marshal(s.Players)
	if err != nil {
		return err
	}

	conn.Cmd("HMSET", key, keyPlayerCount, len(s.Players), keyPlayerList, string(all[:]))

	s.Printf("playerlist: %s\n", string(all[:]))
	s.Printf("playerlist length: %d\n", len(s.Players))

	return nil
}

// OnTick is the simulation update loop
func (s *Simulation) OnTick() {
	conn, err := s.Get()
	if err != nil {
		s.sendErr(err)
		return
	}
	defer s.Put(conn)

	key := makeKey(keySimulationInstance, s.Label)

    onTimeout := func() {
        s.SetPrefix("[UPDATE][ON_TIMEOUT] ")
        s.Printf("timer expired: %s\n", s.Label)
        defer s.SetPrefixDefault()
    }()

    onTick := func() {
        s.SetPrefix("[UPDATE][ON_TICK] ")
        s.Printf("tick: %s\n", s.Label)
        defer s.SetPrefixDefault()

        conn.Cmd("HINCRBY", key, keyFrame, 1)
    }()

    onComplete := func() {
        s.SetPrefix("[UPDATE][ON_DONE] ")
        s.Printf("loop terminated: %s\n", s.Label)
        defer s.SetPrefixDefault()

        s.Timer.Stop()
        s.Ticker.Stop()
    }()

	for {
		select {
		case <-s.Ticker.C:
            onTick()
			// s.SetPrefix("[UPDATE][ON_TICK] ")
			// s.Printf("tick: %s\n", s.Label)
			// s.SetPrefixDefault()

			// conn.Cmd("HINCRBY", key, keyFrame, 1)
		case <-s.Timer.C:
            onTimeout()
			// s.SetPrefix("[UPDATE][ON_TIMEOUT] ")
			// s.Printf("timer expired: %s\n", s.Label)
			// s.SetPrefixDefault()

            return
		case <-s.OnComplete:
            onComplete()
			// s.SetPrefix("[UPDATE][ON_DONE] ")
			// s.Printf("loop terminated: %s\n", s.Label)
			// s.SetPrefixDefault()

			// s.Timer.Stop()
			// s.Ticker.Stop()

			return
		}
	}
}

// OnDestroy terminates the update loop and thus the associated goroutine, and simulation
func (s *Simulation) OnDestroy() {
	s.OnComplete <- true
	close(s.OnComplete)
}

func (s *Simulation) sendErr(err error) {
	s.OnError <- err
	// close(s.OnError)
}

func makeKey(prefix, id string) string {
	return fmt.Sprintf("%s:%s", prefix, id)
}

// GenerateSeed returns a new simulation game world seed
func GenerateSeed() int64 {
	return time.Now().UTC().UnixNano()
}

// GenerateSeedDebug returns the same world seed every time, for debug only
func GenerateSeedDebug() int64 {
	return 1482284596187742126
}
