package session

import (
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet-srv-go/debug"
)

// key
const (
	hk_generatedIDs = "keys:generated"
)

// hash fields
const (
	hf_clientHandleID   = "client_handle_current_key"
	hf_sessionHandleKey = "session_handle_current_key"
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

	p.Cmd("HMSET", hk_generatedIDs, hf_clientHandleID, z, hf_sessionHandleKey, z)

	return kgen, nil
}

// GenerateNext returns an int to be used as a key for a session handle
func (kgen *KeyGenerator) GenerateNextClientID() (*int, error) {
	defer func() {
		kgen.SetPrefixDefault()
	}()

	kgen.SetPrefix("[SESSION][KEY_GEN][CLIENT_HANDLE] ")

	n, err := kgen.Cmd("HINCRBY", hk_generatedIDs, hf_clientHandleID, 1).Int()
	if err != nil {
		return nil, err
	}

	kgen.printGeneratedKey(n)

	return &n, nil
}

func (kgen *KeyGenerator) GenerateNextSessionKey() (*int, error) {
	defer func() {
		kgen.SetPrefixDefault()
	}()

	kgen.SetPrefix("[SESSION][KEY_GEN][SESSION_HANDLE] ")

	n, err := kgen.Cmd("HINCRBY", hk_generatedIDs, hf_sessionHandleKey, 1).Int()
	if err != nil {
		return nil, err
	}

	kgen.printGeneratedKey(n)

	return &n, nil
}

func (kgen *KeyGenerator) printGeneratedKey(k int) {
	kgen.Printf("generated new [id: %d]", k)
}
