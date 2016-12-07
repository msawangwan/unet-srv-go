package main

import (
	"encoding/json"
	"fmt"
	"github.com/msawangwan/unitywebservice/model"
	//"github.com/msawangwan/unitywebservice/data"
)

var JSON = `{
	"hashedgamestate": 89127839172389123
}`

func main() {
	var ss model.PlayerSaveState
	err := json.Unmarshal([]byte(JSON), &ss)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ss)
	var pp []model.PlayerProfile
	pp, _ = model.GetPlayersTest()
	for _, v := range pp {
		fmt.Println(v.Name)
	}
}
