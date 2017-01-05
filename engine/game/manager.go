package game

import (
	"sync"
)

type ActiveTable map[string]*Update

type Manager struct {
	*ActiveTable

	ActiveCount  int
	MaxInstances int

	sync.Mutex
}

func NewGameManager(maxInstances int) *Manager {
	return &Manager{
		ActiveTable:  make(map[string]*Update),
		ActiveCount:  0,
		MaxInstances: maxInstances,
	}
}

func (m *Manager) Add(key string, update *Update) (int, error) {
	m.Lock()
	{
		if m.ActiveCount > m.MaxInstances {
			return -1, errors.New("manager at max capacity")
		} else {
			_, contains := m.ActiveTable[key]
			if !contains {
				m.ActiveTable[key] = update
				m.ActiveCount += 1
			} else {
				return -1, errors.New("instance already in table!")
			}
		}
	}
	m.Unlock()

	return m.ActiveCount, nil
}

func (m *Manager) Access(key string) (*Update, error) {
	return lookUp(key, false)
}

func (n *Manager) Remove(key string) (*Update, error) {
	return lookUp(key, true)
}

func (m *Manager) lookUp(key string, del bool) (*Update, error) {
	var (
		update *Update
		err    error
	)

	m.Lock()
	{
		value, contains := m.ActiveTable[key]
		if contains {
			update = value
			err = nil
			if del {
				delete(m.ActiveTable, key)
				m.ActiveCount -= 1
			}
		} else {
			update = nil
			err = errors.New("no such instance")
		}
	}
	m.Unlock()

	return update, err
}
