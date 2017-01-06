package session

import (
	"strconv"
	"time"

	"github.com/msawangwan/unet/env"
)

const (
	kSessionID          = "session_id"
	kSessionSeed        = "session_seed"
	kSessionPlayerCount = "session_player_count"
)

// Instance is an abstraction for a game session. A game session can be identified
// by a SessionID (unique), a session Seed is used to generate the world
// associated with the session and PlayerCount tracks the number of players
// currently in the session.
type Instance struct {
	SessionID   string `json:"sessionID"`
	Seed        int64  `json:"seed"`
	PlayerCount int    `json:"playerCount"`
}

// Create takes a string and creates a key from it. The key is then cached
// using a redis list which stores the sessionIDs of all currently active
// sessions. It then returns a new session Instance struct.
func Create(e *env.Global, sessionID string) (*Instance, error) {
	k := e.FetchKey_AllActiveSessions()

	res := e.Cmd("SADD", k, sessionID)
	if res.Err != nil {
		return nil, res.Err
	}

	defer func() {
		e.SetPrefix_Debug()
	}()

	e.SetPrefix("[SESSION][CREATE] ")
	e.Printf("created a new session: %s\n", sessionID)

	return &Instance{
			SessionID: sessionID,
			Seed:      generateSeedDebug(),
		},
		nil
}

// Join takes a string, which is used to identify a session instance, and then
// returns that session Instance, increasing the session PlayerCount. Todo:
// handle PlayerCount.
func Join(e *env.Global, gamename string) (*Instance, error) {
	conn, err := e.Get()
	if err != nil {
		return nil, err
	}

	defer func() {
		e.Put(conn)
		e.SetPrefix_Debug()
	}()

	e.SetPrefix("[SESSION][JOIN] ")
	e.Printf("attempting to join game: %s\n", gamename)

	k := e.CreateHashKey_Session(gamename)

	res, err := conn.Cmd("HGETALL", k).Map() // TODO: should probably use a redis WATCHER to prevent race
	if err != nil {
		return nil, err
	}

	var (
		instance *Instance = &Instance{}
	)

	instance.SessionID = res[kSessionID]

	instance.Seed, err = strconv.ParseInt(res[kSessionSeed], 10, 64)
	if err != nil {
		return nil, err
	}

	instance.PlayerCount, err = strconv.Atoi(res[kSessionPlayerCount])
	if err != nil {
		return nil, err
	}

	e.Lock()
	{
		if instance.PlayerCount >= 2 {
			e.Printf("player count greater than 2\n") // TODO: handle this
		} else {
			instance.PlayerCount += 1
			conn.Cmd("HINCRBY", k, kSessionPlayerCount, 1)
		}
	}
	e.Unlock()

	e.Printf("joined game, number of players is: %d\n", instance.PlayerCount)

	return instance, nil
}

// Connect takes an ip (string) and adds it to a list of session connections
func (i *Instance) Connect(e *env.Global, ip string) (bool, *string, error) {
	conn, err := e.Get()
	if err != nil {
		return false, nil, err
	}

	defer func() {
		e.Put(conn)
		e.SetPrefix_Debug()
	}()

	e.SetPrefix("[SESSION][CONNECT] ")

	k := e.CreateListKey_SessionConn(e.CreateHashKey_Session(i.SessionID))

	e.Printf("key for connections: %s\n", k)
	e.Printf("%s connecting to %s\n", ip, i.SessionID)

	res, err := conn.Cmd("LRANGE", k, 0, -1).List()
	if err != nil {
		return false, nil, err
	}

	for _, v := range res {
		if v == ip {
			e.Printf("%s is already connected to the session\n", ip) // TODO: actually handle this, ie return instead of break
			break
		}
	}

	conn.Cmd("RPUSH", k, ip)

	return true, &k, nil
}

// LoadSessionInstanceIntoMemory will add all the properties of a session
// Instance struct into a redis store, which can later be accessed by
func (i *Instance) LoadSessionInstanceIntoMemory(e *env.Global) (*string, error) {
	conn, err := e.Get()
	if err != nil {
		return nil, err
	}

	defer func() {
		e.Put(conn)
		e.SetPrefix_Debug()
	}()

	e.SetPrefix("[SESSION][MAKE_ACTIVE] ")

	if err = conn.Cmd("MULTI").Err; err != nil { // start transaction
		return nil, err
	}

	k := e.CreateHashKey_Session(i.SessionID)

	conn.Cmd("HSET", k, kSessionID, i.SessionID)
	conn.Cmd("HSET", k, kSessionSeed, i.Seed)
	conn.Cmd("HSET", k, kSessionPlayerCount, i.PlayerCount)

	if err = conn.Cmd("EXEC").Err; err != nil {
		return nil, err
	}

	e.Printf("session loaded into memory ...")

	return &k, nil
}

func (i *Instance) KeyFromInstance(e *env.Global) (*string, error) {
	conn, err := e.Get()
	if err != nil {
		return nil, err
	}

	defer func() {
		e.Put(conn)
		e.SetPrefixDefault()
	}()

	e.SetPrefix("[SESSION][KEY_FROM_INSTANCE] ")

	k := e.CreateHashKey_Session(i.SessionID)

	//	if err = conn.Cmd("MULTI").Err; err  != nil { // TODO: use watch
	//		return nil, err
	//	}

	count, err := conn.Cmd("HGET", k, kSessionPlayerCount).Int()
	if err != nil {
		return nil, err
	}

	if count >= 2 {
		e.Printf("session at max capacity (2): %d\n", count)
	}

	count += 1

	err = conn.Cmd("HSET", k, kSessionPlayerCount, count).Err
	if err != nil {
		return nil, err // TODO: rollback??
	}

	return &k, nil
}

func generateSeed() int64 {
	return time.Now().UTC().UnixNano()
}

func generateSeedDebug() int64 {
	return 1482284596187742126
}

type Key struct {
	BareFormat  string `json:"bareFormat"`
	RedisFormat string `json:"redisFormat"`
}
