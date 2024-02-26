package ai

import (
	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
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

	act.Dest = gamestateObj.GetBall().GetPosition()
	act.Dest.SetVec(0, 50)
	act.Dest.SetVec(1, 0)
	act.Dest.SetVec(2, 0)
	act.Dribble = true

	actionList = append(actionList, act)

	return actionList
}
