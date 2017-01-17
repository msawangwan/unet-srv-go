package model

// type Gamestate contains:
// a profile struct and a starmap struct
type GameState struct {
	CurrentProfile        *Profile       `json:"currentProfile"`
	CurrentGameParameters *GameParameter `json:"currentGameParameters"`
}

func NewGameState(p *Profile, gp *GameParameter) *GameState {
	return &GameState{
		CurrentProfile:        p,
		CurrentGameParameters: gp,
	}
}
