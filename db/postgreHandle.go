package db

import (
	"fmt"
	//	"log"

	"database/sql"

	_ "github.com/lib/pq"
)

const (
	kDB_DRIVER   = "postgres"
	kDB_USER     = "postgres"
	kDB_PASSWORD = "1234"
	kDB_DATABASE = "unitywebservice"
)

// type PostgreHandle wraps a db connection pool into postgres db
type PostgreHandle struct {
	*sql.DB
}

// NewPostgreHandle returns an instance of the connection pool, or nil and an
// error if there was one
func NewPostgreHandle() (*PostgreHandle, error) {
	var (
		pg  *PostgreHandle
		db  *sql.DB
		err error
	)

	connstr := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable",
		kDB_USER,
		kDB_PASSWORD,
		kDB_DATABASE,
	)

	db, err = sql.Open(kDB_DRIVER, connstr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	pg = &PostgreHandle{
		DB: db,
	}

	return pg, nil
}
