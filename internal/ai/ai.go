package ai

import (
	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
)

type Ai struct {
	team      gamestate.Team
	ugglan    *PreCalculator
	strutsen  *StrategyFinder
	fiskmasen *RoleExecutor
}

// Constructor for the ai, initializes the client
// and the different components used in the decision pipeline
func NewAi(team gamestate.Team) *Ai {
	ai := &Ai{
		team:      team,
		ugglan:    NewPreCalculator(9000, 6000, 100), // Field length, field width, zone size
		strutsen:  NewPlayFinder(),
		fiskmasen: NewRoleExecutor(),
	}
	return ai
}

// Decides on new actions for the robots
func (ai *Ai) CreateActions(gamestate *gamestate.GameState) ([]action.Action, float64) {

	gameAnalysis := ai.ugglan.Analyse(gamestate)
	score := ai.strutsen.FindStrategy(gamestate, gameAnalysis)
	actions := ai.fiskmasen.GetActions(gamestate, gameAnalysis)
	return actions, score

}
