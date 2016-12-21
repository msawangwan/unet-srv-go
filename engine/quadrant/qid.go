package quadrant

import (
	"sync"
)

// an idGenerator encapsulates all information for tracking quadrant ids
type idCache struct {
	next     int
	assigned map[int]bool
	sync.Mutex
}

// nextID returns the next valid id, and should be run in a goroutine
func (idc *idCache) nextID() int {
	idc.Lock()
	defer idc.Unlock()

increment:
	idc.next++

	if idc.assigned[idc.next] { // in reality: should never be true
		goto increment
	} else {
		idc.assigned[idc.next] = true
	}

	return idc.next
}
