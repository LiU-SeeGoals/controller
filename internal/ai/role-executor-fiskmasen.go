package ai

import (
	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
)

type RoleExecutor struct {
	tornseglaren *PathPlanner
}

func NewRoleExecutor() *RoleExecutor {
	re := &RoleExecutor{
		tornseglaren: NewPathPlanner(),
	}
	return re
}

func (re *RoleExecutor) GetActions(roles *Roles, gs *gamestate.GameState) []action.Action {

	var actionList []action.Action

	act := &action.MoveTo{}
	id := 4

	robot := gs.GetRobot(id, gamestate.Yellow)
	act.Pos = robot.GetPosition()
	act.Id = robot.GetID()

	act.Dest = gs.GetBall().GetPosition()
	act.Dest.SetVec(0, 50)
	act.Dest.SetVec(1, 0)
	act.Dest.SetVec(2, 0)
	act.Dribble = true

	actionList = append(actionList, act)

	return actionList
}
