package game

import (
	"errors"
	"time"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet/debug"
)

func CreateNew(label string, key string, manager *Manager, conns *pool.Pool, log *debug.Log) (*Update, error) {
	var (
		loop *Update
		err  error
	)

	defer func() {
		log.SetPrefixDefault()
	}()

	log.SetPrefix("[GAME ROUTINE][START] ")

	if manager.atMaxCapacity() {
		return nil, errors.New("max game instances running")
	}

	loop, err = NewUpdateRoutine(label, key, conns, log)
	if err != nil {
		return nil, err
	}

	_, err = manager.Add(label, loop)
	if err != nil {
		return nil, err
	}

	endAfter_debug := func() {
		time.Sleep(30 * time.Second)

		log.SetPrefix("[GAME ROUTINE END] ")
		log.Printf("*** *** ***")
		log.Printf("terminated routine %s\n", loop.Label)
		log.Printf("*** *** ***")
		log.SetPrefixDefault()

		loop.OnDestroy()
	}

	if err = loop.Enter("some_player_name"); err != nil {
		return nil, err
	}

	go loop.OnTick()
	go endAfter_debug()

	return loop, nil
}

func EnterExisting(label string, manager *Manager, log *debug.Log) (*Update, error) {
	var (
		loop *Update
		err  error
	)

	defer func() {
		log.SetPrefixDefault()
	}()

	log.SetPrefix("[GAME ROUTINE][ENTER EXISTING] ")

	loop, err = manager.Access(label)
	if err != nil {
		return nil, err
	} else {
		if err = loop.Enter("another_player_name"); err != nil {
			return nil, err
		}

		log.Printf("*** *** ***")
		log.Printf("joined routine %s\n", loop.Label)
		log.Printf("*** *** ***")

		return loop, nil
	}
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

	loop.OnDestroy() // close the loop and clean up ...

	return loop, nil
}
