package prng

import (
	"errors"
	"fmt"
	"sync"
)

// Manager wraps and manages a prng instances
type Manager struct {
	WorldInstance  *Instance         // use for things like world generation
	PlayerInstance map[int]*Instance // by player index
	sync.Mutex
}

func NewInstanceManager(playercount int, s int64) (*Manager, error) {
	var instances = make(map[int]*Instance)
	for i := 1; i < playercount+1; i++ {
		instances[i] = New(s)
	}
	m := &Manager{
		WorldInstance:  New(s),
		PlayerInstance: instances,
	}
	return m, nil
}

func (m *Manager) NextWorldn(max int) int {
	return m.WorldInstance.Intn(max)
}

func (m *Manager) NextWorldf() float32 {
	return m.WorldInstance.Float32()
}

func (m *Manager) NextWorldInRange(min, max float32) float32 {
	return m.WorldInstance.InRange(min, max)
}

func (m *Manager) GetByID(pid int) (*Instance, error) {
	m.Lock()
	defer m.Unlock()

	i, ok := m.PlayerInstance[pid]
	if ok {
		return i, nil
	}
	return nil, errors.New(fmt.Sprintf("prng instance doesn't exist for [player index: %d]", pid))
}
