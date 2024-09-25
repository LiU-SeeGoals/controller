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
		ugglan:    NewPreCalculator(9000, 6000, 100, team), // Field length, field width, zone size
		strutsen:  NewPlayFinder(),
		fiskmasen: NewRoleExecutor(),
	}
	return ai
}

// Decides on new actions for the robots
func (ai *Ai) CreateActions(gamestate *gamestate.GameState) ([]action.Action, float64, float64) {

	// Code by jakob. Not working and should probably not be here,
	// but as a temporary merge it is here now :)
	// If you know where this snipped belong, please edit :)

	// heightMaps := []height_map.HeightMap{
	// 	height_map.HeightMapEnemyGauss{},
	// }
	// // Call FindLowestHeight to get the best position
	// bestX, bestY := height_map.FindLowestHeight(0, float32(5), 36, heightMaps, ai.gamestateObj)
	// // fmt.Println(bestX/1000, bestY/1000)

	// ai.sim_control.TeleportRobot(bestX, bestY, 0, simulation.Team_YELLOW)

	// actions := ai.GenerateMoveActions([]int{0}, []struct{ x, y float64 }{{x: float64(bestX), y: float64(bestY)}})
	// // actions := ai.GenerateMoveActions([]int{0}, []struct{ x, y float64 }{{x: float64(1000), y: float64(1000)}})
	// ai.client.SendActions(actions) // Send actions

	gameAnalysis := ai.ugglan.Analyse(gamestate)
	score, antScore := ai.strutsen.FindStrategy(gamestate, gameAnalysis)
	actions := ai.fiskmasen.GetActions(gamestate, gameAnalysis)
	return actions, score, antScore
}
