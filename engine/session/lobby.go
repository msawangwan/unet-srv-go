package session

import (
	"github.com/msawangwan/unet/db"
	"github.com/msawangwan/unet/env"
)

type Lobby struct {
	Listing []string `json:"listing"`
}

func (l *Lobby) PopulateLobbyList(e *env.Global) error {
	defer func() {
		e.SetPrefix_Debug()
	}()

	e.SetPrefix("[LOBBY INFO] ")
	e.Printf("fetching lobby list ...\n")

	k1 := e.FetchKey_AllActiveSessions()

	r := e.Cmd(db.CMD_SMEMBERS, k1)
	if r.Err != nil {
		return r.Err
	} else {
		sessions, _ := r.List()
		l.Listing = make([]string, len(sessions))
		for i, session := range sessions {
			e.Printf("active session: %s\n", session)
			l.Listing[i] = session
		}

		return nil
	}
}

type LobbyAvailability struct {
	IsAvailable bool `json:"isAvailable"`
}

func (la *LobbyAvailability) CheckAvailability(e *env.Global, name string) error {
	k1 := e.FetchKey_AllActiveSessions()

	res, err := e.Cmd(db.CMD_SISMEMBER, k1, name).Int()
	if err != nil {
		return err
	} else {
		if res == 0 {
			la.IsAvailable = true
		} else {
			la.IsAvailable = false
		}
	}

	return nil
}
