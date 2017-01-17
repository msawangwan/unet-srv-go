package session

import (
	"fmt"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet-srv-go/debug"
)

// GetLobby queries redis and composes an array of strings representing games
func GetLobby(p *pool.Pool, l *debug.Log) ([]string, error) {
	games, err := p.Cmd("SMEMBERS", gameNames()).List()
	if err != nil {
		return nil, err
	}

	listing := make([]string, len(games))

	defer l.PrefixReset()
	l.Prefix("lobby", "populategamelist")
	l.Printf("fetching ...")

	for i, g := range games {
		l.Printf("%d) %s\n", i, g)
		listing[i] = g
	}

	return listing, nil
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
	defer l.PrefixReset()
	l.Prefix("lobby", "postgame")
	l.Printf("adding [game id: %d] [game name: %s] to lobby list", gameid, gamename)

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
