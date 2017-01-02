package db

import (
	"fmt"

	"github.com/mediocregopher/radix.v2/pool"
)

const (
	kCONN_TYPE = "tcp"
	kCONN_SIZE = 10
	kLADDR     = "localhost:6379"
)

// type RedisHandle wraps around a pool of conns for the redis db
type RedisHandle struct {
	*pool.Pool
}

// NewRedisHandle returns a new redis connection handler or nil and an error if
// there was one
func NewRedisHandle() (*RedisHandle, error) {
	var (
		redis *RedisHandle
		conns *pool.Pool
		err   error
	)

	conns, err = pool.New(kCONN_TYPE, kLADDR, kCONN_SIZE)
	if err != nil {
		return nil, err
	}

	redis = &RedisHandle{
		Pool: conns,
	}

	return redis, nil
}

// CreateKey_IsWorldInMemory creates a key from a profile name that is used to
// see if the players world is already loaded in memory
func (rh *RedisHandle) CreateKey_IsWorldInMemory(key string) string {
	return fmt.Sprintf("%s%s", KEY_IS_LOADED_IN_MEMORY, key)
}

// CreateKey_ValidWorldNodes creates a key from a profile name that is used to
// access a hashtable that stores valid world nodes
func (rh *RedisHandle) CreateKey_ValidWorldNodes(key string) string {
	return fmt.Sprintf("%s%s", KEY_WORLD_NODES, key)
}

// CreateHashKey_SessionKey takes a string (game/session name) and from it
// creates and returns a redis hash key (session info).
func (rh *RedisHandle) CreateHashKey_Session(key string) string {
	return fmt.Sprintf("%s:%s", "session", key)
}

// CreateListKey_SessionConn takes a string (a session key) and from it creates
// and returns a redis list key (key for a linked list of conns).
func (rh *RedisHandle) CreateListKey_SessionConn(key string) string {
	return fmt.Sprintf("%s:%s", key, "conn")
}

func (rh *RedisHandle) CreateKey_GameInstance(key string) string {
	return fmt.Sprintf("%s:%s", "game:update", key)
}

func (rh *RedisHandle) FetchKey_AllActiveSessions() string {
	return fmt.Sprintf("%s", "session:all:active")
}
