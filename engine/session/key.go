package session

import (
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet-srv-go/debug"
)

// key
const (
	kSessionKey = "session:key"
)

// KeyGenerator generates the next available int key to be used as a session id
type KeyGenerator struct {
	Next int
	*debug.Log
}

// Key is probably deprecated on the server side
type Key struct {
	BareFormat  string `json:"bareFormat"`
	RedisFormat string `json:"redisFormat"`
}

// NewKeyGenerator is a factory function, reurns an instance of KeyGenerator
func NewKeyGenerator(p *pool.Pool, l *debug.Log) (*KeyGenerator, error) {
	kgen := &KeyGenerator{
		Next: -1,
		Log:  l,
	}

	p.Cmd("SET", kSessionKey, kgen.Next) // TODO: do we care about this error

	return kgen, nil
}

// GenerateNext returns an int to be used as a key for a session handle
func (kgen *KeyGenerator) GenerateNext(p *pool.Pool) (*int, error) {
	conn, err := p.Get()
	if err != nil {
		return nil, err
	}

	defer func() {
		p.Put(conn)
		kgen.SetPrefixDefault()
	}()

	kgen.SetPrefix("[SESSION][KEY_GEN] ")

	n, err := conn.Cmd("INCR", kSessionKey).Int()
	if err != nil {
		return nil, err
	} else {
		kgen.Next = n
		kgen.Printf("registered new session with key (session_id): %d", kgen.Next)
	}

	return &n, nil
}
