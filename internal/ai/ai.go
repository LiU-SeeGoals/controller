package ai

import (
	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"github.com/LiU-SeeGoals/controller/internal/webserver"
	"gonum.org/v1/gonum/mat"
)

type Ai struct {
	gamestateObj *gamestate.GameState
	client       *client.SimClient // TODO change this
	ugglan       *PreCalculator
	strutsen     *StrategyFinder
	hackspetten  *RoleAssigner
	fiskmasen    *RoleExecutor
}

// Constructor for the ai, initializes the client
// and the different components used in the decision pipeline
func NewAi(addr string, gamestateObj *gamestate.GameState) *Ai {
	ai := &Ai{
		ugglan:      NewPreCalculator(9, 6), // TODO: these are field dimensions, should be read from config
		strutsen:    NewPlayFinder(),
		hackspetten: NewRoleAssigner(),
		fiskmasen:   NewRoleExecutor(),

		gamestateObj: gamestateObj,
		client:       client.NewSimClient(addr),
	}
	ai.client.Init()
	return ai
}

// Returns one actions for each robot
func (ai *Ai) decisionPipeline() []action.Action {
	// --- decision pipeline ---
	gameAnalysis := ai.ugglan.Analyse(ai.gamestateObj)
	plays := ai.strutsen.FindStrategy(gameAnalysis)
	roles := ai.hackspetten.AssignRoles(plays)
	actions := ai.fiskmasen.GetActions(roles, ai.gamestateObj)
	return actions
}

// Check if there are any new actions from the webserver
// returns an empty list if no actions were sent
func (ai *Ai) manualControl() []action.Action {
	// --- Manual control ---
	incomming := webserver.GetIncoming() // List of incoming actions
	manualActions := []action.Action{}
	if len(incomming) > 0 {
		// maybe do something with the incoming actions?
	}
	return manualActions
}

// Decides on new actions for the robots, then send them out with the use of the client
// and broadcast the gamestate through the webserver to Gameviewer
func (ai *Ai) CreateAndSendActions() {
	calculatedActions := ai.decisionPipeline() // Calculate new actions

	manualActions := ai.manualControl() // Manual control

	// TODO: Replace calculated actions with the manual ones
	// for the relevant robots. Automatic control should probably
	// be disabled for some time for a robot that has received
	actions := calculatedActions
	if len(manualActions) > 0 {
		actions = manualActions
	}

	actions = ai.GenerateMoveActions([]int{0, 1}, []struct{ x, y float64 }{{x: 0.0, y: 0.0}, {x: 0.0, y: 0.0}})
	ai.client.SendActions(actions)                         // Send actions
	webserver.BroadcastGameState(ai.gamestateObj.ToJson()) // NOTE temporary, will soon change to proto messages

}

// This function can be used to test the MoveTo action.
// To use it simply call it with a slice of ids for the robots
// of interest and the corresponding destinations.
// This will return a list of actions that can be sent to the client.
func (ai *Ai) GenerateMoveActions(robotIDs []int, destinations []struct{ x, y float64 }) []action.Action {
	var actionList []action.Action

	for i, robotID := range robotIDs {
		if i >= len(destinations) {
			break // Prevent out-of-range errors if there are more IDs than destinations
		}

		act := &action.MoveTo{}
		robot := ai.gamestateObj.GetRobot(robotID, gamestate.Yellow)
		act.Pos = robot.GetPosition()
		act.Id = robot.GetID()

		destX := destinations[i].x
		destY := destinations[i].y
		act.Dest = mat.NewVecDense(3, []float64{destX, destY, 0})

		act.Dribble = true // Assuming all moves require dribbling

		actionList = append(actionList, act)
	}

	return actionList
}
