package session

import (
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet/debug"
)

type Lobby struct {
	Listing []string `json:"listing"`
}

func (lobby *Lobby) PopulateLobbyList(p *pool.Pool, l *debug.Log) error {
	defer func() {
		l.SetPrefixDefault()
	}()

	l.SetPrefix("[LOBBY][INFO] ")
	l.Printf("fetching lobby list ...\n")

	k := kSessionAllActive

	r := p.Cmd("SMEMBERS", k)
	if r.Err != nil {
		return r.Err
	} else {
		sessions, _ := r.List()
		lobby.Listing = make([]string, len(sessions))
		for i, session := range sessions {
			l.Printf("active session: %s\n", session)
			lobby.Listing[i] = session
		}

		return nil
	}
}

type LobbyAvailability struct {
	IsAvailable bool `json:"isAvailable"`
}

func (la *LobbyAvailability) CheckAvailability(name string, p *pool.Pool, l *debug.Log) error {
	k := kSessionAllActive

	res, err := p.Cmd("SISMEMBER", k, name).Int()
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
