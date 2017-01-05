package game

import (
	"time"

	"github.com/msawangwan/unet/debug"
)

const (
	kMaxDuration = 1 * time.Hour // more of a saftey, or for special events, increase or remove later
	kTimeout     = 5 * time.Minute
	kTick        = 2500 * time.Millisecond
)

type Update struct {
	Label string `json:"label"`

	Timer  *time.Timer
	Ticker *time.Ticker
	Done   chan bool

	*debug.Log
}

func NewUpdateRoutine(label string, log *debug.Log) *Update {
	return &Update{
		Label:  label,
		Timer:  time.NewTimer(kMaxDuration),
		Ticker: time.NewTicker(kTick),
		Done:   make(chan bool),
		Log:    log,
	}
}

func (u *Update) OnTick() {
	for {
		select {
		case <-u.Timer.C:
			u.Printf("timer fired: %s\n", u.Label)
		case <-u.Ticker.C:
			u.SetPrefix("[UPDATE][ON_TICK] ")
			u.Printf("ticker fired: %s\n", u.Label)
			u.SetPrefixDefault()
		case <-u.Done:
			u.SetPrefix("[UPDATE][ON_DONE] ")
			u.Printf("loop terminated: %s\n", u.Label)
			u.Timer.Stop()
			u.Ticker.Stop()
			u.SetPrefixDefault()
			return
		}
	}
}

func (u *Update) OnDestroy() {
	u.Done <- true
	close(u.Done)
}
