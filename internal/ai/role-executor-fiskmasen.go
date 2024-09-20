package ai

import (
	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"gonum.org/v1/gonum/mat"
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

func (re *RoleExecutor) GetActions(gs *gamestate.GameState, gameAnalysis *GameAnalysis) []action.Action {

	var actionList []action.Action

	myTeam := gs.GetTeam(gameAnalysis.team)

	for _, robot := range myTeam {

		act := action.MoveTo{}
		act.Pos = robot.GetPosition()
		act.Id = robot.GetID()

		anticipatePosition := robot.GetAnticipatedPosition()
		destX := anticipatePosition.AtVec(0)
		destY := anticipatePosition.AtVec(1)
		act.Dest = mat.NewVecDense(3, []float64{destX, destY, 0})

		act.Dribble = true // Assuming all moves require dribbling
		if destX == act.Pos.AtVec(0) && destY == act.Pos.AtVec(1) {
			continue
		}
		actionList = append(actionList, &act)
	}

	return actionList
}
