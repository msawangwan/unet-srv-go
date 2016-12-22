package quadrant

import (
	"sync"
)

type id int

// a store encapsulates all information for tracking and delegating quadrant ids
type store struct {
	next     id
	assigned map[id]bool
	sync.Mutex
}

func NewIDStore(start int) *store {
	return &store{
		next:     id(start),
		assigned: make(map[id]bool),
	}
}

func (s *store) nextAvailable() id {
	s.Lock()
	defer s.Unlock()

increment:
	s.next++

	if s.assigned[s.next] {
		goto increment
	} else {
		s.assigned[s.next] = true
	}

	return s.next
}
