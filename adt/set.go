package adt

import "sync"

const (
	maxsize = 100000 // TODO: enforce this (currently not used anywhere)
)

type StringSet struct {
	store map[string]struct{}
	sync.Mutex
	//sync.RWMutex
}

func NewStringSet() *StringSet {
	return &StringSet{
		store: make(map[string]struct{}),
	}
}

func (ss *StringSet) Sadd(k string) bool {
	if ss.IsMember(k) {
		return false
	}

	ss.Lock()
	defer ss.Unlock()

	ss.store[k] = struct{}{}

	return true
}

func (ss *StringSet) IsMember(k string) bool {
	ss.Lock()
	defer ss.Unlock()

	if _, ex := ss.store[k]; !ex {
		return false
	}

	return true
}

func (ss *StringSet) ListMem() string {
	ss.Lock()
	defer ss.Unlock()

	var s string

	for k, _ := range ss.store {
		s += "\n" + k
	}

	return s
}
