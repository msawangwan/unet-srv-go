package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	//	"log"
	"github.com/msawangwan/unet/util"
)

const (
	kDB_DRIVER   = "postgres"
	kDB_USER     = "postgres"
	kDB_PASSWORD = "1234"
	kDB_DATABASE = "unitywebservice"
)

type postgreManager struct {
	DB *sql.DB
}

var Postgres *postgreManager

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
		util.Log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		util.Log.Fatal(err)
	}

	Postgres = &postgreManager{
		DB: db,
	}

	util.Log.InitMessage("postgres ready ...")
}
