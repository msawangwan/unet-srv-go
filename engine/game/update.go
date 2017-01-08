package game

import (
	"os"
	"runtime"
	"syscall"
	"time"

	"os/signal"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet/debug"
)

const (
	kReportInterval = 5000 * time.Millisecond
)

// type UpdateHandler is the central game loop
type UpdateHandler struct {
	*pool.Pool
	*debug.Log

	Error chan error
	kill  chan os.Signal
}

// NewUpdateHandle returns a new instance of an update handler
func NewUpdateHandle(conns *pool.Pool, log *debug.Log) *UpdateHandler {
	return &UpdateHandler{
		Pool: conns,
		Log:  log,
		kill: make(chan os.Signal, 2),
	}
}

// Run is the core game loop, must be run via goroutine
func (uh *UpdateHandler) Run() {
	for {
		select {
		case err := <-uh.Error:
			uh.SetPrefix("[UPDATE][HANDLE][MAIN][ERROR] ")
			uh.Printf("%s\n", err.Error())
			uh.SetPrefixDefault()
		}
	}
}

// Monitor prints and logs game engine stats, must be run via goroutine
func (uh *UpdateHandler) Monitor() {
	var (
		active int
		avail  int
	)

	signal.Notify(uh.kill, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-uh.kill:
			uh.SetPrefix("[UPDATE][HANDLE][MONITOR] ")
			uh.Printf("terminated, running cleanup (use ctrl+c to exit) ...\n")
			uh.SetPrefixDefault()

			signal.Stop(uh.kill)

			return
		default:
			time.Sleep(kReportInterval)

			active = runtime.NumGoroutine()
			avail = uh.Avail()

			uh.SetPrefix("[UPDATE][HANDLE][MONITOR] ")
			uh.Printf("goroutine count: [%d] available redis conns [%d]\n", active, avail)
			uh.SetPrefixDefault()
		}
	}
}
