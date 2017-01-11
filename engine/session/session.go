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

// client session handle prefixes
const (
	prefixClientHandle           = "client:session:handle"
	prefixClientSessionHandleKey = "client:session:handle:key"
)

// client session handle hash keys
const (
	hashKeyClientName = "name_bare"
)

func makeKey(prefix, id string) string {
	return fmt.Sprintf("%s:%s", prefix, id)
}

type ClientHandle struct{}

func RegisterClient(clientName string, kgen *KeyGenerator, conns *pool.Pool, log *debug.Log) (*int, error) {
	var (
		id *int
	)

	id, err := kgen.GenerateNext()
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
	conn.Cmd("HMSET", val, hashKeyClientName, clientName)

	if err = conn.Cmd("EXEC").Err; err != nil {
		return nil, err
	}

	log.Printf("registered [client name: %s] [client handler id: %d]", clientName, *id)

	return id, nil
}
