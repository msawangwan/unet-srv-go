package game

import (
	"fmt"
	"sync"
	"time"

	"encoding/json"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet/debug"
)

const (
	kMaxPlayersAllowed = 10
)

const (
	kSimulationKey = "sim:key"
)

const (
	kSimulationLabel       = "label"
	kSimulationFrame       = "frame"
	kSimulationPlayerCount = "player_count"
	kSimulationPlayerList  = "player_list"
)

const (
	kDebugTimeout   = 1 * time.Hour
	kConnTimeout    = 5 * time.Minute
	kUpdateInterval = 2500 * time.Millisecond
)

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

func NewSimulation(label string, seed int64, conns *pool.Pool, log *debug.Log) (*Simulation, error) {
	s := &Simulation{
		Label: label,
		Seed:  seed,

		Players: make([]string, 0, kMaxPlayersAllowed),

		OnComplete: make(chan bool),
		OnError:    make(chan error),

		Timer:  time.NewTimer(kDebugTimeout),
		Ticker: time.NewTicker(kUpdateInterval),

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

	key := makeKey(kSimulationKey, s.Label)

	check, err := conn.Cmd("EXISTS", key).Int()
	if err != nil {
		return nil, err
	} else if check == 1 {
		s.Printf("already exists: %s ...", key)
	} else {
		s.Printf("created new simulation: [key %s] [label %s] ...", key, s.Label)
		conn.Cmd("HMSET", key, kSimulationLabel, s.Label)
	}

	return s, nil
}

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

	key := makeKey(kSimulationKey, s.Label)

	n := len(s.Players)
	if n >= kMaxPlayersAllowed {
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

	conn.Cmd("HMSET", key, kSimulationPlayerCount, len(s.Players), kSimulationPlayerList, string(all[:]))

	s.Printf("playerlist: %s\n", string(all[:]))
	s.Printf("playerlist length: %d\n", len(s.Players))

	return nil
}

func (s *Simulation) OnTick() {
	conn, err := s.Get()
	if err != nil {
		s.sendErr(err)
		return
	}
	defer s.Put(conn)

	key := makeKey(kSimulationKey, s.Label)

	for {
		select {
		case <-s.Timer.C:
			s.SetPrefix("[UPDATE][ON_TIMEOUT] ")
			s.Printf("timer expired: %s\n", s.Label)
			s.SetPrefixDefault()
		case <-s.Ticker.C:
			s.SetPrefix("[UPDATE][ON_TICK] ")
			s.Printf("tick: %s\n", s.Label)
			s.SetPrefixDefault()

			conn.Cmd("HINCRBY", key, kSimulationFrame, 1)
		case <-s.OnComplete:
			s.SetPrefix("[UPDATE][ON_DONE] ")
			s.Printf("loop terminated: %s\n", s.Label)
			s.SetPrefixDefault()

			s.Timer.Stop()
			s.Ticker.Stop()

			return
		}
	}
}

func (s *Simulation) OnDestroy() {
	s.OnComplete <- true
	close(s.OnComplete)
}

func (s *Simulation) sendErr(err error) {
	s.OnError <- err
	close(s.OnError)
}
func makeKey(prefix, id string) string {
	return fmt.Sprintf("%s:%s", prefix, id)
}

func GenerateSeed() int64 {
	return time.Now().UTC().UnixNano()
}

func GenerateSeedDebug() int64 {
	return 1482284596187742126
}
