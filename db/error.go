package db

import (
	"errors"
)

type DatabaseError struct {
	Error   error
	Message string
	IsFatal bool
}

// redis-specific errors
var (
	ErrCreatingRedisStore = errors.New("db: error creating redis store")
	ErrStoreAlreadyExists = errors.New("db: redis store already exists")
)
