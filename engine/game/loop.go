package game

import (
	"log"
	"time"
)

/*
 * based on techniques outlined in this article:
 * http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html
 */

var UpdateDispatcher chan chan UpdateLoop
var UpdateQueue = make(chan UpdateLoop, 100)

func init() {
	//StartDispatcher(50)
}

func StartDispatcher(maxLoops int) {
	UpdateDispatcher = make(chan chan UpdateLoop, maxLoops)

	for i := 0; i < maxLoops; i++ {
		log.Println("pooling a loop", i+1)
		l := NewLooper(i+1, UpdateDispatcher)
		l.Start()
	}

	go func() {
		for {
			select {
			case loop := <-UpdateQueue:
				go func() {
					looper := <-UpdateDispatcher
					looper <- loop
				}()
			}
		}
	}()
}

type UpdateLoop struct {
	Label string
}

type Looper struct {
	ID          int
	Loop        chan UpdateLoop
	LooperQueue chan chan UpdateLoop
	Kill        chan bool
}

func NewLooper(id int, looperQueue chan chan UpdateLoop) Looper {
	return Looper{
		ID:          id,
		Loop:        make(chan UpdateLoop),
		LooperQueue: looperQueue,
		Kill:        make(chan bool),
	}
}

func (l *Looper) Start() {
	go func() {
		l.LooperQueue <- l.Loop // add itself to the master queue
		for {
			select {
			case loop := <-l.Loop:
				log.Println("my id", l.ID)
				log.Printf("loop %+v\n", loop)
				time.Sleep(3000 * time.Millisecond) // timestep
				UpdateQueue <- loop                 // feedback
			case <-l.Kill:
				return
			}
		}
	}()
}

func (l *Looper) Stop() {
	go func() {
		l.Kill <- true
	}()
}
