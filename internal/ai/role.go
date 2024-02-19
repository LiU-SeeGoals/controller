package ai

import (
	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
)

type Role interface {
	// return the next action
	NextStep() action.Action
	Assign(id int)
	AssignHeuristic(robots [gamestate.TEAM_SIZE]*gamestate.Robot) int
}

// type Role struct {

// 	ActionHandler ActionHandler
// }
