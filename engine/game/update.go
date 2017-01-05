package game

import (
	"os"
	"runtime"
	"syscall"
	"time"

	"os/signal"

	"github.com/msawangwan/unet/debug"
)

const (
	kReportInterval = 5000 * time.Millisecond
)

type UpdateHandler struct {
	*debug.Log
	kill chan os.Signal
}

func NewUpdateHandle(log *debug.Log) *UpdateHandler {
	return &UpdateHandler{
		Log:  log,
		kill: make(chan os.Signal, 2),
	}
}

//func (uh *UpdateHandler) Run() {
//	for {
//
//	}
//}

// Monitor prints and logs game engine stats, must be run via goroutine
func (uh *UpdateHandler) Monitor() {
	var (
		active int
	)

	signal.Notify(uh.kill, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-uh.kill:
			uh.SetPrefix("[UPDATE HANDLE][MONITOR] ")
			uh.Printf("terminated, running cleanup (use ctrl+c to exit) ...\n")
			signal.Stop(uh.kill)
			uh.SetPrefixDefault()
			return
		default:
			time.Sleep(kReportInterval)
			active = runtime.NumGoroutine()
			uh.SetPrefix("[UPDATE HANDLE][MONITOR] ")
			uh.Printf("goroutine count: %d\n", active)
			uh.SetPrefixDefault()
		}
	}
}
