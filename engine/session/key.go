package session

import (
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet-srv-go/debug"
)

// key
const (
	generatedKeys = "keys:generated"
)

// KeyGenerator generates the next available int key to be used as a session id
type KeyGenerator struct {
	Next int
	*pool.Pool
	*debug.Log
}

// NewKeyGenerator is a factory function, reurns an instance of KeyGenerator
func NewKeyGenerator(p *pool.Pool, l *debug.Log) (*KeyGenerator, error) {
	kgen := &KeyGenerator{
		Next: -1,
		Pool: p,
		Log:  l,
	}

	p.Cmd("SET", generatedKeys, kgen.Next) // TODO: do we care about this error

	return kgen, nil
}

// GenerateNext returns an int to be used as a key for a session handle
func (kgen *KeyGenerator) GenerateNext() (*int, error) {
	conn, err := kgen.Get()
	if err != nil {
		return nil, err
	}

	defer func() {
		kgen.Put(conn)
		kgen.SetPrefixDefault()
	}()

	kgen.SetPrefix("[SESSION][KEY_GEN] ")

	n, err := conn.Cmd("INCR", generatedKeys).Int()
	if err != nil {
		return nil, err
	}

	kgen.Next = n
	kgen.Printf("generated new [id: %d]", kgen.Next)

	return &n, nil
}
