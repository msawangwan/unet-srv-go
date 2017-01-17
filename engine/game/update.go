package game

import (
	"errors"
	"os"
	"runtime"
	"syscall"
	"time"

	"os/signal"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet-srv-go/debug"
)

const (
	statusUpdateInterval = 5000 * time.Millisecond
)

var (
	errTerminatedByAdmin = errors.New("update was terminated by admin")
)

// UpdateHandler is the central game loop
type UpdateHandler struct {
	*pool.Pool
	*debug.Log

	Error chan error
	kill  chan os.Signal
}

// NewUpdateHandle returns a new instance of an update handler
func NewUpdateHandle(errc chan error, conns *pool.Pool, log *debug.Log) *UpdateHandler {
	return &UpdateHandler{
		Pool:  conns,
		Log:   log,
		Error: errc,
		kill:  make(chan os.Signal, 2),
	}
}

// Run is the core game loop, must be run via goroutine
func (uh *UpdateHandler) Run() {
	for {
		select {
		case err := <-uh.Error:
			uh.Prefix("update", "run")
			uh.Printf("%s\n", err.Error())
			uh.PrefixReset()
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
			uh.Prefix("update", "monitor")
			uh.Printf("terminated, running cleanup (use ctrl+c to exit) ...\n")
			uh.Printf("WARNING: database has been flushed due to running in debug mode")
			uh.PrefixReset()

			uh.Cmd("FLUSHDB")

			signal.Stop(uh.kill)
			uh.Error <- errTerminatedByAdmin

			return
		default:
			time.Sleep(statusUpdateInterval)

			active = runtime.NumGoroutine()
			avail = uh.Avail()

			uh.Prefix("update", "monitor")
			uh.Printf("goroutine count: [%d] available redis conns [%d]\n", active, avail)
			uh.PrefixReset()
		}
	}
}
