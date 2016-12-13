package main

import (
	"encoding/json"
	"fmt"
	//	"github.com/msawangwan/unitywebservice/db"
	"github.com/msawangwan/unitywebservice/model"
)

var JSON = `{
	"hashedgamestate": 89127839172389123
}`

func main() {
	var player model.Player
	var players []model.Player
	var err error

	err = json.Unmarshal([]byte(JSON), &player)
	if err != nil {
		fmt.Println(err)
		return
	}

	//	fmt.Printf("a single player %s\n", player)

	players, err = model.SelectAllPlayers()
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range players {
		fmt.Printf("pname: %s\n", v.Name)
		fmt.Printf("puuid: %s\n", v.UUID)
		fmt.Printf("hash: %d\n", v.HashedGameState)
		fmt.Printf("ptime: %v\n", v.DateCreated)
		fmt.Printf("psave: %v\n", v.TimeOfLastSave)
	}

	//	fmt.Printf("a uuid: %s\n", db.CreateUUID())
}
