package game

import (
	"errors"
	"time"

	"github.com/mediocregopher/radix.v2/redis"
	"github.com/msawangwan/unet/debug"
)

func CreateNew(label string, manager *Manager, conn *redis.Client, log *debug.Log) (*Update, error) {
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

	loop = NewUpdateRoutine(label, conn, log)
	endAfter_debug := func() {
		time.Sleep(120 * time.Second)
		log.SetPrefix("[GAME ROUTINE END] ")
		log.Printf("terminated routine %s\n", loop.Label)
		log.SetPrefixDefault()
		loop.OnDestroy()
	}

	go loop.OnTick()
	go endAfter_debug()

	return loop, nil
}
