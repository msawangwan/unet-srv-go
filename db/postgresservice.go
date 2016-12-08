package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

const (
	kDB_DRIVER   = "postgres"
	kDB_USER     = "postgres"
	kDB_PASSWORD = "1234"
	kDB_DATABASE = "unitywebservice"
)

//var Service *sql.DB

type postGreService struct {
	DB *sql.DB
	//logger *log.Logger
}

var PGSconn *postGreService

func init() {
	var db *sql.DB
	var err error

	connstr := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable",
		kDB_USER,
		kDB_PASSWORD,
		kDB_DATABASE,
	)

	db, err = sql.Open(kDB_DRIVER, connstr)

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	PGSconn = &postGreService{
		DB: db,
	}

	log.Printf("[db][db.go][init db: success]\n")
}
