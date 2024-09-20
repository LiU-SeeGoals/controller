package ai

import (
	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"gonum.org/v1/gonum/mat"
)

type Ai struct {
	team        gamestate.Team
	ugglan      *PreCalculator
	strutsen    *StrategyFinder
	hackspetten *RoleAssigner
	fiskmasen   *RoleExecutor
}

// Constructor for the ai, initializes the client
// and the different components used in the decision pipeline
func NewAi(team gamestate.Team) *Ai {
	ai := &Ai{
		team:        team,
		ugglan:      NewPreCalculator(9000, 6000, 1000), // Field length, field width, zone size
		strutsen:    NewPlayFinder(),
		hackspetten: NewRoleAssigner(),
		fiskmasen:   NewRoleExecutor(),
	}
	return ai
}

// Decides on new actions for the robots
func (ai *Ai) CreateActions(gamestate *gamestate.GameState) []action.Action {

	gameAnalysis := ai.ugglan.Analyse(gamestate)
	plays := ai.strutsen.FindStrategy(gamestate, gameAnalysis)
	roles := ai.hackspetten.AssignRoles(plays)
	actions := ai.fiskmasen.GetActions(roles, gamestate)

	actions = ai.GenerateMoveActions(gamestate, []int{0, 1}, []struct{ x, y float64 }{{x: 0.0, y: 0.0}, {x: 0.0, y: 0.0}})

	return actions

}

// This function can be used to test the MoveTo action.
// To use it simply call it with a slice of ids for the robots
func (ai *Ai) GenerateMoveActions(gamestate *gamestate.GameState, robotIDs []int, destinations []struct{ x, y float64 }) []action.Action {
	// of interest and the corresponding destinations.
	// This will return a list of actions that can be sent to the client.
	var actionList []action.Action

	for i, robotID := range robotIDs {
		if i >= len(destinations) {
			break // Prevent out-of-range errors if there are more IDs than destinations
		}

		act := action.MoveTo{}
		robot := gamestate.GetRobot(robotID, ai.team)
		act.Pos = robot.GetPosition()
		act.Id = robot.GetID()

		destX := destinations[i].x
		destY := destinations[i].y
		act.Dest = mat.NewVecDense(3, []float64{destX, destY, 0})

		act.Dribble = true // Assuming all moves require dribbling

		actionList = append(actionList, &act)
	}

	return actionList
}
