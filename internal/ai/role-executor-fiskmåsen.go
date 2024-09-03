package ai

import (
	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"gonum.org/v1/gonum/mat"
	"math"
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

func (re *RoleExecutor) relativeAngle(posA, posB *mat.VecDense) float64 {
	diff := mat.NewVecDense(3, nil)
	diff.SubVec(posB, posA)

	goalAngle := math.Atan2(diff.AtVec(1), diff.AtVec(0))
	currAngle := posA.AtVec(2)

	normalizeAngle := func(angle float64) float64 {
		angle = math.Mod(angle+math.Pi, 2*math.Pi)
		if angle < 0 {
			angle += 2 * math.Pi
		}
		return angle - math.Pi
	}

	return normalizeAngle(goalAngle - currAngle)
}
