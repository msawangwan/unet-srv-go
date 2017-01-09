package session

import (
	"errors"
	"sync"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet/debug"

	"github.com/msawangwan/unet/engine/game"
)

// key format;
// [category]:[label]:[info]

// key
const (
	kSessionHandle = "session::handle"
)

// key values
const (
	kSessionID    = "id"
	kSessionOwner = "owner"
)

// type Handle represents a client session, every client is mapped to a handle and the handle contains:
// - the clients ip
// - what game the client is currently connected to, if any
type Handle struct {
	Owner      string           `json:"owner"`
	OwnerIP    string           `json:"ownerIP"`
	Simulation *game.Simulation `json:"simulation"`
}

func NewHandle(ownerIP string) (*Handle, error) {
	h := &Handle{
		Owner:   "NEED_NAME",
		OwnerIP: ownerIP,
	}
	return h, nil
}

func (h *Handle) AttachToSimulation(sim *game.Simulation) error {
	return nil
}

// type HandleManager is responsible for:
// - creating handles
// - managing a handles lifetime
// - storing them in the db
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
	ErrTableLookupFailed       = errors.New("failed lookup")
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
		key := kSessionHandle + ":" + string(id)
		hm.Cmd("HMSET", key, kSessionID, id, kSessionOwner, handle.Owner)
		hm.Table[id] = handle
		hm.Printf("succeeded in adding session [id: %d] ...\n", id)
		return nil
	}
}

func (hm *HandleManager) Get(id int) (*Handle, error) {
	hm.Lock()

	defer func() {
		hm.Unlock()
		hm.SetPrefixDefault()
	}()

	hm.setPrefix()

	if hm.Table[id] == nil {
		hm.Printf("session handle does not exists, lookup failed [id: %d] ...\n", id)
		return nil, ErrTableLookupFailed
	} else {
		hm.Printf("accessing session handle [id: %d] ...", id)
		return hm.Table[id], nil
	}
}

func (hm *HandleManager) setPrefix() {
	hm.SetPrefix("[SESSION][HANDLE_MANAGER] ")
}