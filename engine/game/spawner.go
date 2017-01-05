package game

import (
	"errors"
	"time"

	"github.com/msawangwan/unet/debug"
)

func CreateNew(manager *Manager, label string, log *debug.Log) (*Update, error) {
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

	loop = NewUpdateRoutine(label, log)
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
