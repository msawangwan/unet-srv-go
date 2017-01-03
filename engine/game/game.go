package game

import (
	"errors"
	"log"
	"runtime"
	"time"

	//	"github.com/msawangwan/unet/db"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/msawangwan/unet/env"
)

type Update struct {
	SessionKey string `json:"sessionKey"`
	Tick       int64  `json:"tick"`
}

func NewInstance(e *env.Global, sessionKey string) (*Update, error) {
	if len(sessionKey) == 0 {
		return nil, errors.New("failed to create an instance of update, invalid id (empty)")
	}

	conn, err := e.Get()
	if err != nil {
		return nil, err
	}

	defer e.Put(conn)

	if err = conn.Cmd("MULTI").Err; err != nil {
		return nil, err
	}

	k := e.CreateHashKey_SessionGameUpdateLoop(sessionKey)

	conn.Cmd("HSET", k, 0, sessionKey)
	conn.Cmd("HSET", k, 1, 0)

	if err = conn.Cmd("EXEC").Err; err != nil {
		return nil, err
	}

	e.Sessions.Add(k, make(chan bool))

	return &Update{
		SessionKey: sessionKey,
		Tick:       0,
	}, nil
}

func (up *Update) Start(e *env.Global) {
	e.Printf("NEW LOOP STARTING %s\n", up.SessionKey)

	k := e.CreateHashKey_SessionGameUpdateLoop(up.SessionKey)

	kill := e.Sessions.Get(k)
	if kill == nil {
		e.Printf("got a nil chan\n")
	}

	conn, err := e.Get()
	if err != nil {
		e.Printf("%s", err.Error())
	}
	defer e.Put(conn)
	e.Add(1)
	go func(c chan bool, r *redis.Client) {
		log.Printf("start game loop for [%s]\n", up.SessionKey)
		for {
			log.Printf("top of loop for [%s]\n", up.SessionKey)
			select {
			case shouldExit := <-c:
				if shouldExit {
					log.Printf("exit\n")
					break
				}
				log.Printf("dont exit\n")
			default:
				log.Printf("ticked\n")
				tick := r.Cmd("HINCRBY", k, 1, 1)
				if tick.Err != nil {
					log.Printf("%s\n", tick.Err)
				}
				time.Sleep(3000 * time.Millisecond)
			}
			runtime.Gosched()
		}
		log.Printf("exited for loop\n")
		//e.Put(conn)
		e.Done()
	}(kill, conn)

	e.Wait()

	e.Printf("returned\n")

	//return nil
}

func (up *Update) End(e *env.Global, sessionKey string) error {
	k := e.CreateHashKey_SessionGameUpdateLoop(sessionKey)

	e.Sessions.Get(k) <- true
	close(e.Sessions.Get(k))

	return nil
}
