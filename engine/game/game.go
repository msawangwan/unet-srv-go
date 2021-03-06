package game

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet-srv-go/debug"
)

type RedisKey string

// the main key -- converts a game id into a redis string
func GameLookupString(gameid int) string {
	return fmt.Sprintf("%s:%d", "game", gameid)
}

// a hash mapping of game param -> value
func GameParamString(gamelookupstr string) string {
	return fmt.Sprintf("%s:%s", gamelookupstr, "param")
}

// a set, lists the players currently in the game
func GamePlayerListString(gamelookupstr string) string {
	return fmt.Sprintf("%s:%s", gamelookupstr, "playerlist")
}

// a key formatted as example: "game:0:player:1"
func GamePlayerNString(gamelookupstr string, index int) string {
	return fmt.Sprintf("%s:player:%d", gamelookupstr, index)
}

// LoadNew loads a new game given a gamename and gameid
func CreateNewGame(gamename string, gameid int, p *pool.Pool, l *debug.Log) (*string, error) {
	gamelookupstr := GameLookupString(gameid)
	gameparamstr := GameParamString(gamelookupstr)

	conn, err := p.Get()
	if err != nil {
		return nil, err
	}

	defer func() {
		p.Put(conn)
		l.PrefixReset()
	}()

	existing, err := conn.Cmd("EXISTS", gamelookupstr).Int()
	if err != nil {
		return nil, err
	} else if existing == 1 {
		return nil, errors.New("fuck you game exists u shit")
	}

	seed := GenerateSeedDebug()

	hasstartedstr := strconv.FormatBool(false)
	seedstring := strconv.FormatInt(seed, 10)
	idstring := strconv.Itoa(gameid)

	l.Prefix("game", "createnew")
	l.Printf("loading a new game [gamename: %s][lookup key: %s][seed: %s] ...", gamename, gameparamstr, seedstring)

	conn.Cmd("HMSET", gamelookupstr, "game_id", idstring, "game_name", gamename, "game_turn_number", -1, "game_player_to_act", 1)
	conn.Cmd("HMSET", gameparamstr, "game_id", idstring, "world_seed", seedstring, "has_started", hasstartedstr)

	return &gamelookupstr, nil
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

type PlayerParameters struct {
	ID    int `json:"id"`
	Index int `json:"index"`
}

func Join(gameid int, playerid int, playername string, p *pool.Pool, l *debug.Log) (*PlayerParameters, error) {
	gamelookupstr := GameLookupString(gameid)
	gameplayerliststring := GamePlayerListString(gamelookupstr)

	conn, err := p.Get()
	if err != nil {
		return nil, err
	}

	defer func() {
		p.Put(conn)
		l.ClearLabel()
	}()

	ismem, err := conn.Cmd("SISMEMBER", gameplayerliststring, playername).Int()
	if err != nil {
		return nil, err
	} else if ismem == 1 {
		return nil, errors.New("player already in-game")
	}

	conn.Cmd("SADD", gameplayerliststring, playername)

	index, err := conn.Cmd("SCARD", gameplayerliststring).Int()
	if err != nil {
		return nil, err
	}

	gameplayerstring := GamePlayerNString(gamelookupstr, index)
	conn.Cmd("HMSET", gameplayerstring, "player_name", playername, "player_id", playerid, "player_index", index)

	allplayers, err := conn.Cmd("SMEMBERS", gameplayerliststring).List()
	if err != nil {
		return nil, err
	}

	l.Label(3, "game", "join")
	l.Printf("printing player list ...")

	for i, player := range allplayers {
		l.Printf("%d) %s\n", i, player)
	}

	return &PlayerParameters{index, playerid}, nil
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
