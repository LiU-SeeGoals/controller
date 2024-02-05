package ai

import (
	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
)

const TEAM_SIZE = 6

type Role interface {
	// return the next action
	NextStep() action.Action
	Assign(id int)
	AssignHeuristic(robots [TEAM_SIZE]*gamestate.Robot) int
}

// type Role struct {

// 	ActionHandler ActionHandler
// }
