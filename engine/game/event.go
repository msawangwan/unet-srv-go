package game

import (
//	"fmt"
)

type Action interface {
	Invoke()
}

type OnJoin struct {
	PlayerName string
}

type OnTurn struct {
	TurnNumber  int
	PlayerToAct int
}
