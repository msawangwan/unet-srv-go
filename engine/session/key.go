package session

import (
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet-srv-go/debug"
)

// key
const (
	keyIdGenerator = "key:generator"
)

// hash fields
const (
	currentClientHandleID   = "client_handlers:curr"
	currentSessionHandleKey = "session_handlers:curr"
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

	p.Cmd("HMSET", keyIdGenerator, currentClientHandleID, z, currentSessionHandleKey, z)

	return kgen, nil
}

// GenerateNext returns an int to be used as a key for a session handle
func (kgen *KeyGenerator) GenerateNextClientID() (*int, error) {
	conn, err := kgen.Get()
	if err != nil {
		return nil, err
	}

	defer func() {
		kgen.Put(conn)
		kgen.SetPrefixDefault()
	}()

	kgen.SetPrefix("[SESSION][KEY_GEN][CLIENT_HANDLE] ")

	n, err := conn.Cmd("HINCRBY", keyIdGenerator, currentClientHandleID, 1).Int()
	if err != nil {
		return nil, err
	}

	kgen.printGeneratedKey(n)

	return &n, nil
}

func (kgen *KeyGenerator) GenerateNextSessionKey() (*int, error) {
	conn, err := kgen.Get()
	if err != nil {
		return nil, err
	}

	defer func() {
		kgen.Put(conn)
		kgen.SetPrefixDefault()
	}()

	kgen.SetPrefix("[SESSION][KEY_GEN][SESSION_HANDLE] ")

	n, err := conn.Cmd("HINCRBY", keyIdGenerator, currentSessionHandleKey, 1).Int()
	if err != nil {
		return nil, err
	}

	kgen.printGeneratedKey(n)

	return &n, nil
}

func (kgen *KeyGenerator) printGeneratedKey(k int) {
	kgen.Printf("generated new [id: %d]", k)
}
