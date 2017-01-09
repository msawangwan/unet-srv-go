package session

import (
	"strconv"
	"sync"
	"time"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet-srv-go/debug"
)

// redis:  list of active sessions
const (
	// active list key
	kSessionAllActive = "session:all:active"
)

// redis: instance
const (
	// entry key
	kSession = "session:"
	// hash keys
	kSessionID          = "session_id"
	kSessionSeed        = "session_seed"
	kSessionPlayerCount = "session_player_count"
)

// type Instance is a session instances abstraction
type Instance struct {
	SessionID   string `json:"sessionID"`
	Seed        int64  `json:"seed"`
	PlayerCount int    `json:"playerCount"`

	sync.Mutex `json:"-"`
}

// Create adds a session to a list of active sessions
func Create(sid string, p *pool.Pool, l *debug.Log) (*Instance, error) {
	k := kSessionAllActive

	res := p.Cmd("SADD", k, sid)
	if res.Err != nil {
		return nil, res.Err
	}

	defer func() {
		l.SetPrefix_Debug()
	}()

	l.SetPrefix("[SESSION][CREATE] ")
	l.Printf("created a new session: %s\n", sid)

	return &Instance{
		SessionID: sid,
		Seed:      generateSeedDebug(),
	}, nil
}

// Join takes a string, which is used to identify a session instance, and then
// returns that session Instance, increasing the session PlayerCount. Todo:
// handle PlayerCount.
func Join(gamename string, p *pool.Pool, l *debug.Log) (*Instance, error) {
	conn, err := p.Get()
	if err != nil {
		return nil, err
	}

	defer func() {
		p.Put(conn)
		l.SetPrefix_Debug()
	}()

	l.SetPrefix("[SESSION][JOIN] ")
	l.Printf("attempting to join game: %s\n", gamename)

	k := kSession + gamename

	res, err := conn.Cmd("HGETALL", k).Map() // TODO: might need WATCH
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

	instance.Lock()
	{
		if instance.PlayerCount >= 2 {
			l.Printf("player count greater than 2\n") // TODO: handle this
		} else {
			instance.PlayerCount += 1
			conn.Cmd("HINCRBY", k, kSessionPlayerCount, 1)
		}
	}
	instance.Unlock()

	l.Printf("joined game, number of players is: %d\n", instance.PlayerCount)

	return instance, nil
}

// Connect takes an ip (string) and adds it to a list of session connections
func (i *Instance) Connect(ip string, p *pool.Pool, l *debug.Log) (bool, *string, error) {
	conn, err := p.Get()
	if err != nil {
		return false, nil, err
	}

	defer func() {
		p.Put(conn)
		l.SetPrefixDefault()
	}()

	l.SetPrefix("[SESSION][CONNECT] ")

	k := "session:" + i.SessionID + ":conn"

	l.Printf("key for connections: %s\n", k)
	l.Printf("%s connecting to %s\n", ip, i.SessionID)

	res, err := conn.Cmd("LRANGE", k, 0, -1).List()
	if err != nil {
		return false, nil, err
	}

	for _, v := range res {
		if v == ip {
			l.Printf("%s is already connected to the session\n", ip) // TODO: actually handle this, ie return instead of break
			break
		}
	}

	conn.Cmd("RPUSH", k, ip)

	return true, &k, nil
}

// LoadSessionInstanceIntoMemory will add all the properties of a session
// Instance struct into a redis store, which can later be accessed by
func (i *Instance) LoadSessionInstanceIntoMemory(p *pool.Pool, l *debug.Log) (*string, error) {
	conn, err := p.Get()
	if err != nil {
		return nil, err
	}

	defer func() {
		p.Put(conn)
		l.SetPrefixDefault()
	}()

	l.SetPrefix("[SESSION][MAKE_ACTIVE] ")

	if err = conn.Cmd("MULTI").Err; err != nil { // start transaction
		return nil, err
	}

	k := kSession + i.SessionID

	conn.Cmd("HSET", k, kSessionID, i.SessionID)
	conn.Cmd("HSET", k, kSessionSeed, i.Seed)
	conn.Cmd("HSET", k, kSessionPlayerCount, i.PlayerCount)

	if err = conn.Cmd("EXEC").Err; err != nil {
		return nil, err
	}

	l.Printf("session loaded into memory ...")

	return &k, nil
}

// func (i *Instance) KeyFromInstance(p *pool.Pool, l *debug.Log) (*string, error) {
// 	conn, err := p.Get()
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer func() {
// 		p.Put(conn)
// 		l.SetPrefixDefault()
// 	}()

// 	l.SetPrefix("[SESSION][KEY_FROM_INSTANCE] ")

// 	k := kSession + i.SessionID

// 	// TODO: use WATCH
// 	count, err := conn.Cmd("HGET", k, kSessionPlayerCount).Int()
// 	if err != nil {
// 		return nil, err
// 	}

// 	if count >= 2 {
// 		l.Printf("session at max capacity (2): %d\n", count)
// 	}

// 	count += 1

// 	err = conn.Cmd("HSET", k, kSessionPlayerCount, count).Err
// 	if err != nil {
// 		return nil, err // TODO: rollback??
// 	}

// 	return &k, nil
// }

func generateSeed() int64 {
	return time.Now().UTC().UnixNano()
}

func generateSeedDebug() int64 {
	return 1482284596187742126
}
