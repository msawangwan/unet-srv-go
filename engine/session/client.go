package session

import (
	"fmt"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet-srv-go/debug"
)

// stores a sorted set of client handles, each member is mapped it a redis hash
const (
	keyAllConnectedClientHandles = "client:session:handle:all"
)

// client session handle key prefixes
const (
	prefixClientHandle           = "client:session:handle"
	prefixClientSessionHandleKey = "client:session:handle:key"
)

// client session handle hash fields
const (
	hfClientName = "name_bare"
	hfAttachedSessionHandle
)

func makeKey(prefix, id string) string {
	return fmt.Sprintf("%s:%s", prefix, id)
}

// ClientHandle is currently not in use
type ClientHandle struct{}

// RegisterClient returns an id mapped to the passed in client (player) name, this id is used to uniquely identify the client on future requests
func RegisterClient(clientName string, kgen *KeyGenerator, conns *pool.Pool, log *debug.Log) (*int, error) {
	var (
		id *int
	)

	id, err := kgen.GenerateNextClientID()
	if err != nil {
		return nil, err
	}

	conn, err := conns.Get()
	if err != nil {
		return nil, err
	}

	defer func() {
		conns.Put(conn)
		log.SetPrefixDefault()
	}()

	log.SetPrefix("[SESSION][REGISTER_CLIENT] ")

	key := makeKey(prefixClientHandle, string(*id))
	val := makeKey(prefixClientSessionHandleKey, clientName)

	if err = conn.Cmd("MULTI").Err; err != nil {
		return nil, err
	}

	conn.Cmd("ZADD", key, *id, val)
	conn.Cmd("HMSET", val, hfClientName, clientName)

	if err = conn.Cmd("EXEC").Err; err != nil {
		return nil, err
	}

	log.Printf("registered [client name: %s] [client handler id: %d]", clientName, *id)

	return id, nil
}

// MapToClient maps a session handle to a client in the redis store
func MapToClient(int id, kgen *KeyGenerator, conns *pool.Pool, log *debug.Log) (*int, error) {
	id, err := kgen.GenerateNextSessionKey()
	if err != nil {
		return nil, error
	}

	conn, err = conns.Get()
	if err != nil {
		return nil, error
	}

	defer func() {
		conns.Put(conn)
		log.SetPrefixDefault()
	}()

	log.SetPrefix("[SESSION][MAP_SESSION_CLIENT] ")
	// LEFT OFF HERE

	return nil
}
