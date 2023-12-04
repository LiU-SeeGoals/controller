package ai

import (
	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"gonum.org/v1/gonum/mat"
)

type Ai struct {
	gamstate     *gamestate.GameState
	Grsim_client *client.GrsimClient
}

// Method used for testing actions,
// a proper test should be implemented
func (ai *Ai) TestActions() {
	act := &action.Move{}
	act.Pos = mat.NewVecDense(3, nil)
	act.Pos.SetVec(0, 0)
	act.Pos.SetVec(1, 0)
	act.Pos.SetVec(2, 0) //gamestate.GetRobot(0, gamestate.Yellow)//ai.gamstate.yellow_team[0].pos
	act.Dest = mat.NewVecDense(3, nil)
	act.Dest.SetVec(0, 4)
	act.Dest.SetVec(1, 0)
	act.Dest.SetVec(2, 0)
	act.Dribble = true

	//act := &action.Kick{}
	//act.Kickspeed = 10

	//act := &action.Dribble{}
	//act.Dribble = true

	var action []action.Action
	action = append(action, act)

	ai.Grsim_client.SendActions(action)
}

func (ai *Ai) Update() {
	ai.TestActions()

}

func NewAi(gamestate *gamestate.GameState, addr string) *Ai {
	ai := &Ai{}

	ai.gamstate = gamestate
	ai.Grsim_client = client.NewGrsimClient(addr)
	ai.Grsim_client.Init()
	return ai
}
