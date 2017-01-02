package session

import (
	"strconv"
	"time"

	"github.com/msawangwan/unet/env"
)

type Instance struct {
	SessionID   string `json:"sessionID"`
	Seed        int64  `json:"seed"`
	PlayerCount int    `json:"playerCount"`
}

func Create(e *env.Global, sessionID string) (*Instance, error) {
	k := e.FetchKey_AllActiveSessions()

	res := e.Cmd("SADD", k, sessionID)
	if res.Err != nil {
		return nil, res.Err
	}

	e.Printf("created a new session: %s\n", sessionID)

	return &Instance{
			SessionID: sessionID,
			Seed:      generateSeedDebug(),
		},
		nil
}

func Join(e *env.Global, gamename string) (*Instance, error) {
	e.Printf("attempting to join game: %s\n", gamename)

	k := e.CreateKey_SessionInstance(gamename)

	res, err := e.Cmd("HGETALL", k).Map()
	if err != nil {
		return nil, err
	}

	var (
		instance *Instance = &Instance{}
	)

	instance.SessionID = res["0"]
	instance.Seed, err = strconv.ParseInt(res["1"], 10, 64)
	if err != nil { // TODO: handle this better
		return nil, err
	}
	instance.PlayerCount, err = strconv.ParseInt(res["2"], 10, 32)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (i *Instance) LoadSessionInstanceIntoMemory(e *env.Global) error {
	conn, err := e.Get()
	if err != nil {
		return err
	}

	defer func() {
		e.Put(conn)
		e.SetPrefix_Debug()
	}()

	e.SetPrefix("[LOADING SESSION INTO MEMORY] ")

	if err = conn.Cmd("MULTI").Err; err != nil { // start transaction
		return err
	}

	e.Printf("started tx ...")

	k := e.CreateKey_SessionInstance(i.SessionID)

	conn.Cmd("HSET", k, 0, i.SessionID)
	conn.Cmd("HSET", k, 1, i.Seed)
	conn.Cmd("HSET", k, 2, i.PlayerCount)

	if err = conn.Cmd("EXEC").Err; err != nil {
		return err
	}

	e.Printf("executed tx ...")
	e.Printf("session loaded into memory ...")

	return nil
}

func generateSeed() int64 {
	return time.Now().UTC().UnixNano()
}

func generateSeedDebug() int64 {
	return 1482284596187742126
}
