package session

import (
	"errors"
	"sync"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet/debug"
)

type Handle struct {
	Owner string `json:"player"`
}

func NewHandle(owner string) (*Handle, error) {
	h := &Handle{
		Owner: owner,
	}
	return h, nil
}

type HandleManager struct {
	Table    map[int]*Handle
	capacity int

	*pool.Pool
	*debug.Log

	sync.Mutex
}

func NewHandleManager(capacity int, conns *pool.Pool, log *debug.Log) (*HandleManager, error) {
	hm := &HandleManager{
		Table:    make(map[int]*Handle, capacity),
		capacity: capacity,

		Pool: conns,
		Log:  log,
	}
	return hm, nil
}

var (
	ErrHandleAlreadyRegistered = errors.New("failed to add session instance")
)

func (hm *HandleManager) Add(id int, handle *Handle) error {
	hm.Lock()

	defer func() {
		hm.Unlock()
		hm.SetPrefixDefault()
	}()

	hm.setPrefix()

	if hm.Table[id] != nil {
		hm.Printf("failed to add session [id: %d] ...\n", id)
		return ErrHandleAlreadyRegistered
	} else {
		hm.Table[id] = handle
		hm.Printf("succeeded in adding session [id: %d] ...\n", id)
		return nil
	}
}

func (hm *HandleManager) setPrefix() {
	hm.SetPrefix("[SESSION][HANDLE_MANAGER] ")
}
