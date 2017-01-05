package game

import (
	"errors"
	"time"

	//"github.com/mediocregopher/radix.v2/redis"
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet/debug"
)

func CreateNew(label string, key string, manager *Manager, conns *pool.Pool, log *debug.Log) (*Update, error) {
	var (
		loop *Update
	)

	defer func() {
		log.SetPrefixDefault()
	}()

	log.SetPrefix("[GAME ROUTINE START] ")

	if manager.atMaxCapacity() {
		return nil, errors.New("max game instances running")
	}

	loop = NewUpdateRoutine(label, key, conns, log)
	manager.Add(label, loop)

	endAfter_debug := func() {
		time.Sleep(30 * time.Second)

		log.SetPrefix("[GAME ROUTINE END] ")
		log.Printf("*** *** ***")
		log.Printf("terminated routine %s\n", loop.Label)
		log.Printf("*** *** ***")
		log.SetPrefixDefault()

		loop.OnDestroy()
	}

	go loop.OnTick()
	go endAfter_debug()

	return loop, nil
}

func EndActive(label string, key string, manager *Manager, log *debug.Log) (*Update, error) {
	loop, err := manager.Remove(label)
	if err != nil {
		return nil, err
	}

	log.SetPrefix("[GAME ROUTINE END] ")
	log.Printf("*** *** ***")
	log.Printf("terminated routine %s\n", loop.Label)
	log.Printf("*** *** ***")
	log.SetPrefixDefault()

	loop.OnDestroy() // close the loop and call clean up functions ...

	return loop, nil
}
