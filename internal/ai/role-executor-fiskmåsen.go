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

func (re *RoleExecutor) GetActions(roles *[]Role, gamestate *gamestate.GameState) []action.Action {

	var actionList []action.Action

	// gigantic switch case in the works
	for _, role := range *roles {
		actionList = append(actionList, role.NextStep())

		// switch role.role {
		// case Keeper:
		// 	actionList = append(actionList, Keeper(role.robotId, gamestate))
		// default:
		// 	fmt.Println("Panic")
		// }
	}

	return actionList
}
