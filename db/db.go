package data

import (
	"github.com/sausheong/gwp/Chapter_2_Go_ChitChat/chitchat/data"
	"log"
)

var PostGreService *sql.DB

func init() {
	var err error
	PostGreService, err = sql.Open("postgres", "dbname=unitywebservice sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if err = PostGreService.Ping(); err != nil {
		log.Fatal(err)
	}
	return
}
