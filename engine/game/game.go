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

// LoadNew loads a new game given a clientID
func LoadNew(cid int, id int, p *pool.Pool, l *debug.Log) error {
	v := fmt.Sprintf("%s:%d", phk_gameHandleSimulation, id)

	l.Prefix("game", "loadnew")
	l.PrefixReset()
	l.Printf("loading a new game [gameid: %d][lookup key: %s] created by [client: %d] ...", id, v, cid)

	p.Cmd("HMSET", hk_gameHandleKey, strconv.Itoa(id), v)

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
