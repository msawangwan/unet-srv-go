package manager

import (
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet-srv-go/debug"
)

const (
	hk_idDispenser = "keys:generated"
)

// hash fields
const (
	hf_clientHandleID   = "client_handle_current_key"
	hf_sessionHandleKey = "session_handle_current_key"
	hf_gameHandleID     = "game_handle_current_key"
	hf_playerID         = "player_id_current_key"
)

// KeyGenerator retrieves assignable keys from the redis store
type KeyGenerator struct {
	*pool.Pool
	*debug.Log
}

// NewKeyGenerator is a factory function, reurns an instance of KeyGenerator
func NewKeyGenerator(p *pool.Pool, l *debug.Log) (*KeyGenerator, error) {
	kgen := &KeyGenerator{
		Pool: p,
		Log:  l,
	}

	z := -1

	p.Cmd("HMSET", hk_idDispenser, hf_clientHandleID, z, hf_sessionHandleKey, z, hf_gameHandleID, z)

	return kgen, nil
}

// GenerateNext returns an int to be used as a key for a session handle
func (kgen *KeyGenerator) GenerateNextClientID() (*int, error) {
	kgen.Prefix("session", "keygen", "clientid")
	defer kgen.PrefixReset()

	n, err := kgen.Cmd("HINCRBY", hk_idDispenser, hf_clientHandleID, 1).Int()
	if err != nil {
		return nil, err
	}

	kgen.printGeneratedKey(n)

	return &n, nil
}

func (kgen *KeyGenerator) GenerateNextSessionKey() (*int, error) {
	kgen.Prefix("session", "keygen", "sessionkey")
	defer kgen.PrefixReset()

	n, err := kgen.Cmd("HINCRBY", hk_idDispenser, hf_sessionHandleKey, 1).Int()
	if err != nil {
		return nil, err
	}

	kgen.printGeneratedKey(n)

	return &n, nil
}

func (kgen *KeyGenerator) GenerateNextGameID() (*int, error) {
	kgen.Prefix("game", "keygen", "gameid")
	defer kgen.PrefixReset()

	n, err := kgen.Cmd("HINCRBY", hk_idDispenser, hf_gameHandleID, 1).Int()
	if err != nil {
		return nil, err
	}

	kgen.printGeneratedKey(n)

	return &n, nil
}

func (kgen *KeyGenerator) GenerateNextPlayerID() (*int, error) {
	kgen.Prefix("game", "keygen", "playerid")
	defer kgen.PrefixReset()

	n, err := kgen.Cmd("HINCRBY", hk_idDispenser, hf_playerID, 1).Int()
	if err != nil {
		return nil, err
	}

	kgen.printGeneratedKey(n)

	return &n, nil
}

func (kgen *KeyGenerator) printGeneratedKey(k int) {
	kgen.Printf("generated new [id: %d]", k)
}
