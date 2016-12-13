package model

//import (
//	"github.com/msawangwan/unitywebservice/db"
//)

type ProfileSearch struct {
	Name        string `json: "name"`
	IsAvailable bool   `json: "isAvailable"`
}

//func CheckAvailability(string name) (*ProfileSearch, error) {
//	var ps *ProfileSearch = &ProfileSearch{
//		Name:        name,
//		IsAvailable: false,
//	}

//	conn, err := db.Redis.DB.Get()
//	if err != nil {
//		return ps, err
//	}
//	defer db.Redis.DB.Put(conn) // remember to return the connection
//reply

//	return ps, nil
//}
