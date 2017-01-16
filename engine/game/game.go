package game

import (
	"fmt"
	"strconv"
	"time"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet-srv-go/debug"
)

// hash keys
const (
	// game db sim
	hk_gameHandleKey = "game:handle:table" // maps gameid -> concat(prefix, [gname|id])
)

// hash fields
const (
	// game db fields
	hf_gameKey = "game_key"
	hf_seed    = "game_seed"
)

const (
	phk_gameHandleSimulation = "game:handle:simulation"
)

// type GameHandle is not used currently
type GameHandle struct {
	prefix string `json:"-"`
}

// TODO: add to map, aka some sort of mananger
func NewGameHandle(id int) *GameHandle {
	return &GameHandle{
		prefix: phk_gameHandleSimulation,
	}
}

// LoadNew loads a new game given a gamename and gameid
func CreateNewHandler(gamename string, id int, p *pool.Pool, l *debug.Log) error {

	v := fmt.Sprintf("%s:%s", phk_gameHandleSimulation, gamename)

	conn, err := p.Get()
	if err != nil {
		return err
	}

	defer func() {
		p.Put(conn)
		l.PrefixReset()
	}()

	l.Prefix("game", "createnewhandler")

	seed := GenerateSeedDebug()
	// seed := GenerateSeed()
	seedstring := strconv.FormatInt(seed, 10)
	idstring := strconv.Itoa(id)

	l.Printf("loading a new game [gamename: %s][lookup key: %s] ...", gamename, v)
	l.Printf("seed [%s]", seedstring)

	conn.Cmd("HSET", hk_gameHandleKey, idstring, v)
	conn.Cmd("HMSET", v, hf_gameKey, idstring, hf_seed, seedstring)

	return nil
}

// GenerateSeed returns a new simulation game world seed
func GenerateSeed() int64 {
	return time.Now().UTC().UnixNano()
}

// GenerateSeedDebug returns the same world seed every time, for debug only
func GenerateSeedDebug() int64 {
	return 1482284596187742126
}
