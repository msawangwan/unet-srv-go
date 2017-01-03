package env

import (
	"sync"
)

type SessionTable struct {
	Active      map[string]chan bool
	ActiveCount int
	MaxActive   int

	sync.Mutex
}

func NewSessionTable(maxSessionPerHost int) *SessionTable {
	return &SessionTable{
		Active:      make(map[string]chan bool),
		ActiveCount: -1,
		MaxActive:   maxSessionPerHost,
	}
}

func (st *SessionTable) Add(key string, c chan bool) {
	st.Lock()
	{
		st.Active[key] = c
		st.ActiveCount += 1
	}
	st.Unlock()
}

func (st *SessionTable) Get(key string) chan bool {
	st.Lock()
	{
		c, ok := st.Active[key]
		if ok {
			return c
		}
	}
	st.Unlock()
	return nil
}
