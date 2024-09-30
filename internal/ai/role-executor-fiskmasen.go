package ai

import (
	// "fmt"
	"github.com/LiU-SeeGoals/controller/internal/action"
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

func (re *RoleExecutor) GetActions(gs *state.GameState, gameAnalysis *GameAnalysis) []action.Action {

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
		// fmt.Println("Robot", act.Id, "moving to", destX, destY, "from", act.Pos.AtVec(0), act.Pos.AtVec(1))
		actionList = append(actionList, &act)
	}

	return actionList
}
