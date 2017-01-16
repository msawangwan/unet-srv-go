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
	//	var (
	//		id *int
	//	)

	//id, err := kgen.GenerateNextClientID()
	//	if err != nil {
	//		return nil, err
	//	}

	conn, err := conns.Get()
	if err != nil {
		return err
	}

	defer func() {
		conns.Put(conn)
		//log.SetPrefixDefault()
		log.PrefixReset()
	}()

	//log.SetPrefix("[SESSION][REGISTER_CLIENT] ")

	//clientID := strconv.Itoa(*id)
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

	//return id, nil
	return nil
}

func IsMapped(chid int, conns *pool.Pool, log *debug.Log) (bool, error) {
	ch := strconv.Itoa(chid)

	chkey := conns.Cmd("HGET", hk_allClientHandles, ch)
	if chkey.Err != nil {
		return false, chkey.Err
		//return false, errors.New("error checking client/session mapping redis")
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
		//	log.SetPrefixDefault()
		log.PrefixReset()
	}()

	//log.SetPrefix("[SESSION][MAP_SESSION_CLIENT] ")
	//log.Printf("querying redis store for client with [handle id : %d] ...", id)

	//sid, err := kgen.GenerateNextSessionKey()
	//if err != nil {
	//	return nil, err
	//}

	//log.Printf("mapping a new session [session key: %d] to client [handle id: %d] ...", sid, chid)

	clientID := strconv.Itoa(chid)
	ch, err := conn.Cmd("HGET", hk_allClientHandles, clientID).Str()
	if err != nil {
		return err
	}

	sessionID := strconv.Itoa(sid)
	conn.Cmd("HSET", ch, hf_attachedSessionHandleByID, sessionID)

	log.Prefix("session", "mapclient")
	log.Printf("mapped session [key: %d] to client [id: %d][handler key: %s]", sid, chid, ch)

	return nil
}
