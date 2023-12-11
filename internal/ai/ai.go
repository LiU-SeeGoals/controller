package ai

import (
	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
)

type Ai struct {
	gamestate    *gamestate.GameState
	Grsim_client *client.GrsimClient
}

// Method used for testing actions,
// a proper test should be implemented
func (ai *Ai) TestActions() {

	var actionList []action.Action
	act := &action.Move{}
	id := 4

	robot := ai.gamestate.GetRobot(id, false)
	act.Pos = robot.GetPosition()
	act.Id = robot.GetID()

	act.Dest = ai.gamestate.GetBall().GetPosition()
	act.Dest.SetVec(0, 4)
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

	ai.Grsim_client.SendActions(actionList)
}

func (ai *Ai) Update() {
	ai.TestActions()

}

func NewAi(gamestate *gamestate.GameState, addr string) *Ai {
	ai := &Ai{}

	ai.gamestate = gamestate
	ai.Grsim_client = client.NewGrsimClient(addr)
	ai.Grsim_client.Init()
	return ai
}
