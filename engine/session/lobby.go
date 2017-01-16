package session

import (
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet-srv-go/debug"
)

// hash keys
const (
	hk_gameLobbyList = "game:name:id" // a map of gamename -> started (bool)
)

// set keys
const (
	sk_namesInUse = "game:name:in-use" // set of names that are already taken
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

	r := p.Cmd("SMEMBERS", sk_namesInUse) // TODO: this hasn't been fixed
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

// CheckAvailability checks a given name against a list of all active simulations to check for uniqueness
func IsGameNameUnique(name string, p *pool.Pool) (bool, error) {
	res, err := p.Cmd("SISMEMBER", sk_namesInUse, name).Int()
	if err != nil {
		return false, err
	} else if res == 1 {
		return false, nil
	}

	return true, nil
}

// PostGameToLobby adds the game name to a set of names
func PostGameToLobby(name string, p *pool.Pool) error {
	conn, err := p.Get()
	if err != nil {
		return err
	}

	defer func() {
		p.Put(conn)
	}()

	if err = conn.Cmd("MULTI").Err; err != nil {
		return err
	}

	conn.Cmd("SADD", sk_namesInUse, name)
	conn.Cmd("HMSET", hk_gameLobbyList, name, false)

	if err = conn.Cmd("EXEC").Err; err != nil {
		return err
	}

	return nil
}
