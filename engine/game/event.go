package game

type Action interface {
	Invoke()
}

type OnJoin struct {
	PlayerName string
}
