package game

import (
	"time"
)

// Frame is a type
type Frame struct {
	SessionID string `json:"sessionID"`

	HasStarted bool `json:"hasStarted"`
	HasWinner  bool `json:"hasWinner"`

	ConnectionCount   int `json:"connectionCount"`
	PlayerToAct       int `json:"playerToAct"`
	CurrentTurnNumber int `json:"currentTurnNumber"`
	PacketID          int `json:"packetID"`

	Timestamp time.Time `json:"timeStamp"`
}
