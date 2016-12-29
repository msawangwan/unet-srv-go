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

	err = redis.createNameDB()
	if err != nil {
		if err != ErrStoreAlreadyExists {
			return nil, err
		} else {
			fmt.Printf("%s", err)
		}
	}

	return redis, nil
}

// createNameDB is a setup routine that cretes a store for user names in use
func (rh *RedisHandle) createNameDB() error {
	conn, err := rh.Get()
	if err != nil {
		return err
	}
	defer rh.Put(conn)

	query := conn.Cmd(CMD_EXISTS, KEY_NAMES_TAKEN)
	if query.Err != nil {
		return query.Err
	} else {
		result, _ := query.Int()
		if result == 1 {
			return ErrStoreAlreadyExists
		}
	}

	if err = conn.Cmd(CMD_SADD, KEY_NAMES_TAKEN, VAL_INIT).Err; err != nil { // TODO: verify this is valid syntax
		return err
	} else {
		return nil
	}
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
