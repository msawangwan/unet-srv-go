package session

import (
	"fmt"
	"strconv"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"

	"github.com/msawangwan/unet-srv-go/debug"
)

// prefixes to forge keys from
const (
	phk_clientHandleByID = "client:handle"
)

// hash keys
const (
	hk_allClientHandles = "client:handle:all" // hash : maps id -> client:handle:id
)

// hash fields
const (
	hf_clientName                = "name_bare"
	hf_attachedSessionHandleByID = "attached_session_handle_id"
)

func marshallIntPointer(np *int) string {
	return strconv.Itoa(*np)
}

// wrapper to print the key we made
func forge(prefix string, val string, l *debug.Log) string {
	k := fmt.Sprintf("%s:%s", prefix, val)
	l.Printf("forged [key: %s]", k)
	return k
}

// ClientHandle is currently not in use
type ClientHandle struct{}

// RegisterClient returns an id mapped to the passed in client (player) name, this id is used to uniquely identify the client on future requests
func RegisterClient(clientName string, clientID int, conns *pool.Pool, log *debug.Log) error {
	conn, err := conns.Get()
	if err != nil {
		return err
	}

	defer func() {
		conns.Put(conn)
		log.PrefixReset()
	}()

	cid := strconv.Itoa(clientID)
	key := forge(phk_clientHandleByID, cid, log) // format: [client:handle:id]

	if err = conn.Cmd("MULTI").Err; err != nil {
		return err
	}

	conn.Cmd("HSET", hk_allClientHandles, cid, key)  // register id -> client:handle (stores the forged key for later retrival)
	conn.Cmd("HSET", key, hf_clientName, clientName) // set a field that stores the current client:handles in-game name

	if err = conn.Cmd("EXEC").Err; err != nil {
		return err
	}

	log.Prefix("session", "registerclient")
	log.Printf("registered [client name: %s] [client handler id: %d]", clientName, clientID)

	return nil
}

func IsMapped(chid int, conns *pool.Pool, log *debug.Log) (bool, error) {
	ch := strconv.Itoa(chid)

	chkey := conns.Cmd("HGET", hk_allClientHandles, ch)
	if chkey.Err != nil {
		return false, chkey.Err
	} else if chkey.IsType(redis.Nil) {
		return false, nil
	}

	return true, nil
}

// MapToClient maps a session handle to a client in the redis store, then
// returns the newly generated session handle id
func MapToClient(chid int, sid int, conns *pool.Pool, log *debug.Log) error {
	conn, err := conns.Get()
	if err != nil {
		return err
	}

	defer func() {
		conns.Put(conn)
		log.PrefixReset()
	}()

	clientID := strconv.Itoa(chid)
	sessionKey := strconv.Itoa(sid)

	ch, err := conn.Cmd("HGET", hk_allClientHandles, clientID).Str()
	if err != nil {
		return err
	}

	conn.Cmd("HSET", ch, hf_attachedSessionHandleByID, sessionKey)

	log.Prefix("session", "mapclient")
	log.Printf("mapped session [key: %d] to client [id: %d][handler key: %s]", sid, chid, ch)

	return nil
}
