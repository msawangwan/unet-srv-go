package db

import (
	"errors"
)

// redis-specific errors
var (
	ErrCreatingRedisStore = errors.New("db: error creating redis store")
	ErrStoreAlreadyExists = errors.New("db: redis store already exists")
)
