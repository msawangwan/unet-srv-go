package session

import (
	"errors"
	"sync"

	"github.com/msawangwan/unet/debug"
)

const (
	kMaxPlayersAllowed = 10 // should be 2 but 10 for debug purposes
)

type Manager struct {
	Instances map[string]*Active
	sync.Mutex
	*debug.Log
}

func NewManager(log *debug.Log) *Manager {
	return &Manager{
		Instances: make(map[string]*Active),
		Log:       log,
	}
}

type Active struct {
	*Instance
	Players []string
}

func (m *Manager) Add(key string, instance *Instance) (int, error) {
	m.SetPrefix("[SESSION][MANAGER][ADD] ")

	defer func() {
		m.SetPrefixDefault()
	}()

	m.Lock()
	{
		m.Printf("attempting to add session: %s\n", key)
		_, contains := m.Instances[key]
		if !contains {
			m.Printf("success adding session: %s\n")
			m.Instances[key] = &Active{
				Instance: instance,
				Players:  make([]string, 0, kMaxPlayersAllowed),
			}
		} else {
			m.Printf("failed to add session: %s\n")
			return len(m.Instances), errors.New("instance already in table!")
		}
	}
	m.Unlock()

	return len(m.Instances), nil
}

func (m *Manager) AddPlayer(key, player string) (*Active, error) {
	a, err := m.lookUp(key, false)
	if err != nil {
		return nil, err
	}

	m.SetPrefix("[SESSION][MANAGER][ADD_PLAYER] ")

	m.Lock()
	{
		n := len(a.Players)
		m.Printf("current player count: %d", n)
		if (n + 1) <= cap(a.Players) {
			m.Printf("added player: %s [%s]", key, player)
			a.Players = append(a.Players, player)
		}
	}
	m.Unlock()

	return a, nil
}

func (m *Manager) AccessAllPlayers(key string) ([]string, error) {
	a, err := m.lookUp(key, false)
	if err != nil {
		return nil, err
	}

	return a.Players, nil
}

func (m *Manager) Access(key string) (*Active, error) {
	return m.lookUp(key, false)
}

func (m *Manager) Remove(key string) (*Active, error) {
	return m.lookUp(key, true)
}

func (m *Manager) lookUp(key string, del bool) (*Active, error) {
	var (
		a   *Active
		err error
	)

	m.Lock()
	{
		value, contains := m.Instances[key]
		if contains {
			a = value
			err = nil
			if del {
				delete(m.Instances, key)
			}
		} else {
			a = nil
			err = errors.New("no such instance")
		}
	}
	m.Unlock()

	return a, err
}
