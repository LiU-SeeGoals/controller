package ai

import (
	"fmt"
	"math"

	"github.com/LiU-SeeGoals/controller/internal/action"
	. "github.com/LiU-SeeGoals/controller/internal/logger"
	"github.com/LiU-SeeGoals/controller/internal/info"
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
	const distance_threshold = 200
	const angle_threshold = 0.1
	distance_achieved := distance_left <= distance_threshold
	angle_diff := math.Abs(float64(curr_pos.Angle - target_position.Angle))
	angle_achieved := angle_diff <= angle_threshold
	return distance_achieved && angle_achieved
}

func (m *MoveToBall) GetID() info.ID {
	return m.id
}

