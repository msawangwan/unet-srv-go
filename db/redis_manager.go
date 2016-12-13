package db

import (
	"errors"
	"github.com/mediocregopher/radix.v2/pool"
	"log"
)

const (
	CONN_TYPE = "tcp"
	CONN_SIZE = 10
	LADDR     = "localhost:6379"

	K_NAMES_NOT_AVAIL = "profile_name:taken"

	VAL_INIT = "_init_"

	CMD_EXIST   = "EXISTS"
	CMD_S       = "SET"
	CMD_G       = "GET"
	CMD_SADDMEM = "SADD"
	CMD_SISMEM  = "SISMEMBER"
	CMD_SETMEM  = "SMEMBERS"
)

type redisManager struct {
	DB *pool.Pool
}

var Redis *redisManager

func init() {
	var db *pool.Pool
	var err error

	db, err = pool.New(CONN_TYPE, LADDR, CONN_SIZE)

	if err != nil {
		log.Panic(err)
	}

	Redis = &redisManager{
		DB: db,
	}

	err = Redis.CreateNameDatabase()

	if err != nil {
		log.Printf("non-fatal error when creating name db -- %v\n", err)
	}
}

var (
	ErrOnInitDB        = errors.New("redis: error on setup")
	ErrAlreadyExistsDB = errors.New("redis: db already exists")
)

func (rm *redisManager) CreateNameDatabase() error {
	conn, err := rm.DB.Get()
	if err != nil {
		return err
	}
	defer rm.DB.Put(conn)

	query := conn.Cmd(CMD_EXIST, K_NAMES_NOT_AVAIL)
	if query.Err != nil {
		return query.Err
	} else {
		result, _ := query.Int()
		if result == 1 {
			return ErrAlreadyExistsDB
		}
	}

	if conn.Cmd(CMD_SADDMEM, K_NAMES_NOT_AVAIL, VAL_INIT).Err != nil {
		return ErrOnInitDB
	} else {
		return nil
	}
}
