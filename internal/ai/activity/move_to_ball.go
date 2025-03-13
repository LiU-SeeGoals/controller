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
	ballPos, err := gi.State.GetBall().GetPosition()
	if err != nil {
		Logger.Errorf("Position retrieval failed - Ball: %v\n", err)
		return NewStop(m.id).GetAction(gi)
	}

	robotPos, err := gi.State.GetRobot(m.id, m.team).GetPosition()
	if err != nil {
		Logger.Errorf("Position retrieval failed - Robot: %v\n", err)
		return NewStop(m.id).GetAction(gi)
	}

	angleToBall := robotPos.AngleToPosition(ballPos)

	target := info.Position{X: ballPos.X, Y: ballPos.Y, Z: 0, Angle: angleToBall}
	move := NewMoveToPosition(m.team, m.id, target)
    fmt.Println("ball: ", target)
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
	const distance_threshold = 90 // WARN: Magic number
	const angle_threshold = 0.1
	distance_achieved := distance_left <= distance_threshold

	angleToBall := curr_pos.AngleToPosition(target_position)
	angle_diff := curr_pos.AngleDistance(info.Position{ X: 0, Y: 0, Z: 0, Angle: angleToBall })
	angle_achieved := angle_diff <= angle_threshold
	return distance_achieved && angle_achieved
}

func (m *MoveToBall) GetID() info.ID {
	return m.id
}

