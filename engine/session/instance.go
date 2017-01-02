package session

import (
	"time"

	"github.com/msawangwan/unet/env"
)

type Instance struct {
	Seed      int64  `json:"seed"`
	SessionID string `json:"sessionID"`
}

func Create(e *env.Global, sessionID string) (*Instance, error) {
	k := e.FetchKey_AllActiveSessions()
	res := e.Cmd("SADD", k, sessionID)
	if res.Err != nil {
		return nil, res.Err
	}

	e.Printf("created a new session: %s\n", sessionID)

	return &Instance{
			Seed:      generateSeedDebug(),
			SessionID: sessionID,
		},
		nil
}

func Join(e *env.Global, gamename string) (*Instance, error) {
	k := e.CreateKey_SessionInstance(gamename)

	res, err := e.Cmd("HGETALL", k).Map()
	if err != nil {
		return nil, err
	}

	var (
		instance *Instance = &Instance{}
	)

	instance.Seed = res["0"]
	instance.SessionID = res["1"]

	return instance, nil
}

func (i *Instance) LoadSessionInstanceIntoMemory(e *env.Global) error {
	if conn, err := e.Get(); err != nil {
		return err
	}

	defer func() {
		e.Put(conn)
		e.SetPrefix_Debug()
	}()

	e.SetPrefix("[load session instance into memory ...]")

	if err = conn.Cmd("MULTI").Err; err != nil { // start transaction
		return err
	}

	e.Printf("started tx ...")

	k := e.CreateKey_SessionInstance(i.SessionID)

	conn.Cmd("HSET", k, 0, i.Seed)
	conn.Cmd("HSET", k, 1, i.SessionID)

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
