package game

import (
	"github.com/msawangwan/unet/debug"
	"log"
	"runtime"
	"sync"
	"time"
)

func GameLoop(label string) {
	log.Printf("started a new loop: %s\n", label)

	tick := time.NewTicker(3000 * time.Millisecond)

	for _ = range tick.C {
		log.Printf("tick %s\n", label)
	}

	//for {
	//	select {
	//	case <-tick.C:
	//		log.Printf("tick %s\n", label)
	//	}
	//}
}

func timeo(d time.Duration) <-chan bool {
	to := make(chan bool, 1)
	go func() {
		time.Sleep(d)
		to <- true
	}()
	return to
}

// type Updater implements a loop to be run in a goroutine
type Updater interface {
	Update(log *debug.Log)
}

type UpdateMonitor struct{}

func (um *UpdateMonitor) Update(l *debug.Log) {
	var (
		activeCount int
	)

	for {
		time.Sleep(5000 * time.Millisecond)
		activeCount = runtime.NumGoroutine()
		l.SetPrefix("[UPDATE MONITOR] ")
		l.Printf("current number of active goroutines: %d\n", activeCount)
		l.SetPrefix_Debug()
	}
}

const (
	Loopers    = 5
	SubLoopers = 3
	SUBLOOPS   = 3
	Tasks      = 20
	SubTasks   = 10
	WORKERS    = 5
	SUBWORKERS = 3
	SUBTASKS   = 3
	TASKS      = 20
)

func subworker(subtasks chan int) {
	for {
		task, ok := <-subtasks
		if !ok {
			return
		}
		time.Sleep(time.Duration(task) * time.Millisecond)
		log.Println(task)
	}
}

func worker(tasks <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		task, ok := <-tasks
		if !ok {
			return
		}

		subtasks := make(chan int)
		for i := 0; i < SUBWORKERS; i++ {
			task1 := task * i
			subtasks <- task1
		}
		for i := 0; i < SUBTASKS; i++ {
			task1 := task * i
			subtasks <- task1
		}
		close(subtasks)
	}
}

func amain() {
	var wg sync.WaitGroup
	wg.Add(WORKERS)
	tasks := make(chan int)

	for i := 0; i < WORKERS; i++ {
		go worker(tasks, &wg)
	}

	for i := 0; i < TASKS; i++ {
		tasks <- i
	}

	close(tasks)
	wg.Wait()
}
