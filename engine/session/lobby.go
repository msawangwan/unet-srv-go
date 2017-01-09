package session

import (
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet-srv-go/debug"
)

// key
const (
	keySessionLobbyList = "session:lobby:all"
)

// Lobby represents a list of attached (to session handles) simulations
type Lobby struct {
	Listing []string `json:"listing"`
}

// PopulateLobbyList hits the redis store to generate the latest lobby list
func (lobby *Lobby) PopulateLobbyList(p *pool.Pool, l *debug.Log) error {
	defer func() {
		l.SetPrefixDefault()
	}()

	l.SetPrefix("[LOBBY][INFO] ")
	l.Printf("fetching lobby list ...\n")

	r := p.Cmd("SMEMBERS", keySessionLobbyList)
	if r.Err != nil {
		return r.Err
	}

	sessions, _ := r.List()
	lobby.Listing = make([]string, len(sessions))

	for i, session := range sessions {
		l.Printf("active session: %s\n", session)
		lobby.Listing[i] = session
	}

	return nil
}

// LobbyAvailability wraps a bool, might be unnecessary
type LobbyAvailability struct {
	IsAvailable bool `json:"isAvailable"`
}

// CheckAvailability checks a given name against a list of all active simulations to check for uniqueness
func (la *LobbyAvailability) CheckAvailability(name string, p *pool.Pool, l *debug.Log) error {
	res, err := p.Cmd("SISMEMBER", keySessionLobbyList, name).Int()
	if err != nil {
		return err
	}

	if res == 0 {
		la.IsAvailable = true
	} else {
		la.IsAvailable = false
	}

	return nil
}
