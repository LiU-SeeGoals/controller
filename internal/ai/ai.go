package ai

import (
	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"gonum.org/v1/gonum/mat"
)

type Ai struct {
	gamestate     *gamestate.GameState
	client        *client.GrsimClient // TODO change this
	preCalculator *PreCalculator
	playFinder    *PlayFinder
	roleAssigner  *RoleAssigner
	roleExecutor  *RoleExecutor
}

func NewAi(gamestate *gamestate.GameState, addr string) *Ai {
	ai := &Ai{
		preCalculator: NewPreCalculator(),
		playFinder:    NewPlayFinder(),
		roleAssigner:  NewRoleAssigner(),
		roleExecutor:  NewRoleExecutor(),

		gamestate: gamestate,
		client:    client.NewGrsimClient(addr),
	}
	ai.client.Init()
	return ai
}

// Method used for testing actions,
// a proper test should be implemented
func (ai *Ai) TestActions() {

	var actionList []action.Action
	act := &action.MoveTo{}
	id := 4

	robot := ai.gamestate.GetRobot(id, gamestate.Blue)
	act.Pos = robot.GetPosition()
	act.Id = robot.GetID()

	act.DestPos = mat.NewVecDense(2, []float64{5000, 0})
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

func (ai *Ai) Update() {
	gameAnalysis := ai.preCalculator.Process(ai.gamestate)
	plays := ai.playFinder.FindPlays(gameAnalysis)
	roles := ai.roleAssigner.AssignRoles(plays)
	actions := ai.roleExecutor.GetActions(roles, ai.gamestate)

	ai.client.SendActions(actions)

	//ai.TestActions()

}
