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
	ballPos, err := gi.State.GetBall().GetEstimatedPosition()
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
	return move.GetAction(gi)
}

func (m *MoveToBall) Achieved(gi *info.GameInfo) bool {
	return gi.State.GetBall().GetPossessor() == gi.State.GetRobot(m.id, m.team)
}

func (m *MoveToBall) GetID() info.ID {
	return m.id
}

