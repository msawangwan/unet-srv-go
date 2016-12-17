package model

type GameState struct {
	CurrentProfile *Profile `json:"currentProfile"`
	CurrentStarMap *StarMap `json:"currentStarMap"`
}

func NewGameState(p *Profile, sm *StarMap) *GameState {
	return &GameState{
		CurrentProfile: p,
		CurrentStarMap: sm,
	}
}
