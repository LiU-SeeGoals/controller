package info

type GameInfo struct {
	State  *GameState
	Status *GameStatus
	Field  *GameField
}

func NewGameInfo(capacity int) *GameInfo {
	return &GameInfo{
		State:  NewGameState(capacity),
		Status: NewGameStatus(),
		Field:  NewGameField(),
	}
}
