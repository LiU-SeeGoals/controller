package ai

import (
	"fmt"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"github.com/LiU-SeeGoals/controller/internal/webserver"
)

type Ai struct {
	gamestate     *gamestate.GameState
	client        *client.SimClient // TODO change this
	preCalculator *PreCalculator
	playFinder    *PlayFinder
	roleAssigner  *RoleAssigner
	roleExecutor  *RoleExecutor
}

func NewAi(addr string, gamestate *gamestate.GameState) *Ai {
	ai := &Ai{
		preCalculator: NewPreCalculator(9, 6),
		playFinder:    NewPlayFinder(),
		roleAssigner:  NewRoleAssigner(),
		roleExecutor:  NewRoleExecutor(),

		gamestate: gamestate,
		client:    client.NewSimClient(addr),
	}
	ai.client.Init()
	return ai
}

func (ai *Ai) handleIncoming(incomming []action.ActionDTO) []action.Action {
	fmt.Println("Received a new action (gamestate)")

	// TODO also set manual control for the robot that is controlled

	// for _, act := range incomming {
	// 	switch act.Action {
	// 	case robot_action.ActionType_MOVE_ACTION:
	// 		pos := mat.NewVecDense(3, []float64{float64(act.PosX), float64(act.PosY), float64(act.PosW)})
	// 		dest := mat.NewVecDense(3, []float64{float64(act.DestX), float64(act.DestY), float64(act.DestW)})
	// 		gs.AddAction(&action.Move{act.Id, pos, dest, act.Dribble})
	// 	case robot_action.ActionType_INIT_ACTION:
	// 		gs.AddAction(&action.Init{act.Id})
	// 	case robot_action.ActionType_ROTATE_ACTION:
	// 		gs.AddAction(&action.Rotate{act.Id, int(act.PosW)})
	// 	case robot_action.ActionType_KICK_ACTION:
	// 		standardKickSpeed := 1
	// 		gs.AddAction(&action.Kick{act.Id, standardKickSpeed})
	// 	case robot_action.ActionType_MOVE_TO_ACTION:
	// 		dest := mat.NewVecDense(3, []float64{float64(act.DestX), float64(act.DestY)})
	// 		gs.AddAction(&action.SetNavigationDirection{act.Id, dest})
	// 	case robot_action.ActionType_STOP_ACTION:
	// 		gs.AddAction(&action.Stop{act.Id})
	// 	}
	// }
	return nil

}

// Method used for testing actions,
// a proper test should be implemented

func (ai *Ai) Update() {
	// --- AI ---
	gameAnalysis := ai.preCalculator.Process(ai.gamestate)
	plays := ai.playFinder.FindPlays(gameAnalysis)
	roles := ai.roleAssigner.AssignRoles(plays)
	actions := ai.roleExecutor.GetActions(roles, ai.gamestate)
	
	// --- Manual control ---
	incomming := webserver.GetIncoming() // List of incoming actions
	manualActions := []action.Action{}
	if len(incomming) > 0 {
		manualActions = ai.handleIncoming(incomming) // If we got new actions --> then handle them
	}
	
	// Replace the manual actions with the calculated actions
	for i := 0; i < len(manualActions); i++ {
		actions[i] = manualActions[i]
	}

	ai.client.SendActions(actions) // Send actions

	//ai.TestActions()

}
func (ai *Ai) TestActions() {

	var actionList []action.Action
	act := &action.MoveTo{}
	id := 4

	robot := ai.gamestate.GetRobot(id, gamestate.Yellow)
	act.Pos = robot.GetPosition()
	act.Id = robot.GetID()

	act.Dest = ai.gamestate.GetBall().GetPosition()
	act.Dest.SetVec(0, 50)
	act.Dest.SetVec(1, 0)
	act.Dest.SetVec(2, 0)
	act.Dribble = true

	actionList = append(actionList, act)

	//for id := 0; id < 6; id++ {
	//
	//	act := &action.Move{}
	//
	//	robot := ai.gamestate.GetRobot(id, false)
	//	act.Pos = robot.GetPosition()
	//	act.Id = robot.GetID()
	//
	//	act.Dest = ai.gamestate.GetBall().GetPosition() //mat.NewVecDense(3, nil)
	//	//act.Dest.SetVec(0, 4)
	//	//act.Dest.SetVec(1, 0)
	//	//act.Dest.SetVec(2, 0)
	//	act.Dribble = true
	//
	//	actionList = append(actionList, act)
	//}

	ai.client.SendActions(actionList)
}
