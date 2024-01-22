package ai

import (
	"fmt"

	"github.com/LiU-SeeGoals/controller/internal/gamestate"
)

type PreCalculator struct {
}

func NewPreCalculator() *PreCalculator {
	pc := &PreCalculator{}
	return pc
}

type Data struct { // name is a placeholder
}

func (pc *PreCalculator) Process(gamestate *gamestate.GameState) *Data {
	fmt.Println("Ugglan")
	return &Data{}
}
