package session

type Connection struct {
	IsConnected bool   `json:"isConnected"`
	Address     string `json:"address"`
	Message     string `json:"message"`
}
