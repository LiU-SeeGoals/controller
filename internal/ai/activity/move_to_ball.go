package ai

import (
	"fmt"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/info"
	. "github.com/LiU-SeeGoals/controller/internal/logger"
)

type MoveToBall struct {
	GenericComposition
	// MovementComposition
}

func (m *MoveToBall) String() string {
	return fmt.Sprintf("(Robot %d, MoveToBall()", m.id)
}

func NewMoveToBall(team info.Team, id info.ID) *MoveToBall {
	return &MoveToBall{
		GenericComposition: GenericComposition{
			team: team,
			id:   id,
		},
	}
}

func (m *MoveToBall) GetAction(gi *info.GameInfo) action.Action {
	target_position, err := gi.State.GetBall().GetPosition()

	if err != nil {
		Logger.Errorf("Position retrieval failed - Ball: %v\n", err)
		return NewStop(m.id).GetAction(gi)
	}

	move := NewMoveToPosition(m.team, m.id, target_position)
	return move.GetAction(gi)
}

func (m *MoveToBall) Achieved(gi *info.GameInfo) bool {
	target_position, err := gi.State.GetBall().GetPosition()

	if err != nil {
		Logger.Errorf("Position retrieval failed - Ball: %v\n", err)
		return false
	}

	curr_pos, err := gi.State.GetTeam(m.team)[m.id].GetPosition()
	if err != nil {
		Logger.Errorf("Position retrieval failed - Robot: %v\n", err)
		return false
	}
	distance_left := curr_pos.Distance(target_position)
	const distance_threshold = 100
	distance_achieved := distance_left <= distance_threshold

	return distance_achieved
}
