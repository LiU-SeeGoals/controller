package info

import "github.com/LiU-SeeGoals/proto_go/ssl_vision"

type GameInfo struct {
	State  *GameState
	Status *GameStatus
	Field  *ssl_vision.SSL_GeometryFieldSize
}

func NewGameInfo(capacity int) *GameInfo {
	return &GameInfo{
		State:  NewGameState(capacity),
		Status: NewGameStatus(),
		Field:  nil,
	}
}
