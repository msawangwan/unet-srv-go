package session

import (
	"fmt"
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet-srv-go/debug"
)

// hash keys
//const (
//	hk_gameLobbyList = "game:lobby-map" // a map of gamename -> started (bool)
//)

// set keys
//const (
//	sk_namesInUse = "game:lobby:names-in-use" // set of names that are already taken
//)

// Lobby represents a list of attached (to session handles) simulations
type Lobby struct {
	Listing []string `json:"listing"`
}

func GetLobby(p *pool.Pool, l *debug.Log) ([]string, error) {
	defer func() {
		l.PrefixReset()
	}()

	l.Prefix("lobby", "populategamelist")
	l.Printf("fetching ...")

	games, err := p.Cmd("SMEMBERS", gameNames()).List()
	if err != nil {
		return nil, err
	}

	listing := make([]string, len(games))

	for i, g := range games {
		l.Printf("%d) %s\n", g)
		listing[i] = g
	}

	return listing, nil
}

// PopulateLobbyList hits the redis store to generate the latest lobby list
func (lobby *Lobby) PopulateLobbyList(p *pool.Pool, l *debug.Log) error {
	defer func() {
		l.PrefixReset()
	}()

	l.Prefix("lobby", "info")
	l.Printf("fetching lobby list ...\n")

	sessions, err := p.Cmd("SMEMBERS", gameNames()).List()
	if err != nil {
		return err
	}

	lobby.Listing = make([]string, len(sessions))

	for i, session := range sessions {
		l.Printf("active session: %s\n", session)
		lobby.Listing[i] = session
	}

	return nil
}

// CheckAvailability checks a given name against a list of all active simulations to check for uniqueness
func IsGameNameUnique(name string, p *pool.Pool) (bool, error) {
	res, err := p.Cmd("SISMEMBER", gameNames(), name).Int()
	if err != nil {
		return false, err
	} else if res == 1 {
		return false, nil
	}

	return true, nil
}

// PostGameToLobby adds the game name to a set of names
func PostGameToLobby(gameid int, gamename string, p *pool.Pool, l *debug.Log) error {
	p.Cmd("SADD", gameNames(), gamename)
	p.Cmd("HSET", gameLobby(), gamename, gameid)

	return nil
}

// set key, for fast member test
func gameNames() string {
	return fmt.Sprintf("%s:%s", "lobby", "names")
}

// hash map (also fast...?) to get key mapped to name
func gameLobby() string {
	return fmt.Sprintf("%s:%s", "lobby", "games")
}
