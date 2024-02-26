package ai

import (
	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"gonum.org/v1/gonum/mat"
)

type RoleExecutor struct {
	pathPlanner *PathPlanner
}

func NewRoleExecutor() *RoleExecutor {
	re := &RoleExecutor{
		pathPlanner: NewPathPlanner(),
	}
	return re
}

func (re *RoleExecutor) GetActions(roles *Roles, gamestateObj *gamestate.GameState) []action.Action {

	var actionList []action.Action

	act := &action.MoveTo{}
	id := 4

	robot := gamestateObj.GetRobot(id, gamestate.Blue)
	act.Pos = robot.GetPosition()
	act.Id = robot.GetID()

	act.DestPos = mat.NewVecDense(2, []float64{5000, 0})
	act.Dribble = true

	actionList = append(actionList, act)

	return actionList
}
