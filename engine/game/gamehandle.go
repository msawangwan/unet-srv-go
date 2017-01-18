package game

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet-srv-go/debug"
)

// hash fields
const (
	hf_gameKey = "game_key"
	hf_seed    = "world_seed"
	hf_players = "game_player_list"
)

func GameHandlerString(gameid int) string {
	return fmt.Sprintf("%s:%d", "game", gameid)
}

// LoadNew loads a new game given a gamename and gameid
func CreateNewGame(gamename string, gameid int, p *pool.Pool, l *debug.Log) (*string, error) {
	gamehandlerstring := GameHandlerString(gameid)

	conn, err := p.Get()
	if err != nil {
		return nil, err
	}

	defer func() {
		p.Put(conn)
		l.PrefixReset()
	}()

	exists, err := conn.Cmd("SISMEMBER", "game:lookup_key", gamehandlerstring).Int()
	if err != nil {
		return nil, err
	} else if exists == 1 {
		return nil, errors.New("fuck you")
	}

	seed := GenerateSeedDebug()

	seedstring := strconv.FormatInt(seed, 10)
	idstring := strconv.Itoa(gameid)

	l.Prefix("game", "createnew")
	l.Printf("loading a new game [gamename: %s][lookup key: %s][seed: %s] ...", gamename, gamehandlerstring, seedstring)

	conn.Cmd("SADD", "game:list", gamehandlerstring)
	conn.Cmd("HMSET", gamehandlerstring, hf_gameKey, idstring, hf_seed, seedstring, hf_players, "")

	return &gamehandlerstring, nil
}

func GetExistingGameByID(gamename string, p *pool.Pool, l *debug.Log) (*int, error) {
	lobbygameid, err := p.Cmd("HGET", gameLobby(), gamename).Str()
	if err != nil {
		return nil, err
	}

	gameid, err := strconv.Atoi(lobbygameid)
	if err != nil {
		return nil, err
	}

	return &gameid, nil
}

func Join(gameid int, playername string, p *pool.Pool, l *debug.Log) error {
	gamehandlerstring := GameHandlerString(gameid)
	gameplayerliststring := fmt.Sprintf("%s:%s", gamehandlerstring, "playerlist")

	conn, err := p.Get()
	if err != nil {
		return err
	}

	defer func() {
		p.Put(conn)
		l.PrefixReset()
	}()

	ismem, err := conn.Cmd("SISMEMBER", gameplayerliststring, playername).Int()
	if err != nil {
		return err
	} else if ismem == 1 {
		return errors.New("player already in-game")
	}

	conn.Cmd("SADD", gameplayerliststring, playername)

	allplayers, err := conn.Cmd("SMEMBERS", gameplayerliststring).List()
	if err != nil {
		return err
	}

	l.Prefix("game", "join")
	l.Printf("printing player list ...")

	for i, player := range allplayers {
		l.Printf("%d) %s\n", i, player)
	}

	return nil
}

// GenerateSeed returns a new simulation game world seed
func GenerateSeed() int64 {
	return time.Now().UTC().UnixNano()
}

// GenerateSeedDebug returns the same world seed every time, for debug only
func GenerateSeedDebug() int64 {
	return 1482284596187742126
}

// set of game names
func gameNames() string {
	return fmt.Sprintf("%s:%s", "lobby", "names")
}

// hash mapped name -> key
func gameLobby() string {
	return fmt.Sprintf("%s:%s", "lobby", "games")
}
