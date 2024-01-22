package ai

import (
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
)

type PreCalculator struct {
}

func NewPreCalculator() *PreCalculator {
	pc := &PreCalculator{}
	return pc
}

type GameAnalysis struct {
}

func (pc *PreCalculator) Process(gamestate *gamestate.GameState) *GameAnalysis {

	return &GameAnalysis{}
}
