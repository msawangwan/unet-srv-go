package main

import (
	"encoding/json"
	"fmt"
	"github.com/msawangwan/unitywebservice/model"
)

var JSON = `{
	"hashedgamestate": 89127839172389123
}`

func main() {
	var ss model.SaveState
	err := json.Unmarshal([]byte(JSON), &ss)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ss)
}
