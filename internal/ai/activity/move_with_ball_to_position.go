package ai

import (
	"fmt"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/info"
	. "github.com/LiU-SeeGoals/controller/internal/logger"
)

type MoveWithBallToPosition struct {
	GenericComposition
	targetPosition info.Position
}

func (m *MoveWithBallToPosition) String() string {
	return fmt.Sprintf("Robot %d, MoveWithBallToPosition(%v)", m.id, m.targetPosition)
}

func NewMoveWithBallToPosition(team info.Team, id info.ID, dest info.Position) *MoveWithBallToPosition {
	return &MoveWithBallToPosition{
		GenericComposition: GenericComposition{
			team: team,
			id:   id,
		},
		targetPosition: dest,
	}
}

func (fb *MoveWithBallToPosition) GetAction(gi *info.GameInfo) action.Action {

	robot := gi.State.GetRobot(fb.id, fb.team)
	robotPosition, err := robot.GetPosition()
	if err != nil {
		Logger.Errorf("Position retrieval failed - Robot: %v\n", err)
		return NewStop(fb.id).GetAction(gi)
	}
	if !robot.IsActive() {
		return NewStop(fb.id).GetAction(gi)
	}

	ball := gi.State.GetBall()

	// If we lost the ball, go get it
	if ball.GetPossessor() != robot { // WARN: Magic number
		Logger.Debug("MoveWithBallToPosition: Lost possession of ball")

		move := NewMoveToBall(fb.team, fb.id)
		return move.GetAction(gi)
	}

	Logger.Debug("MoveWithBallToPosition: Moving with ball to position")
	// Passed all checks, is in possession
	// move with ball to target position
	act := action.MoveTo{
		Id:      int(fb.id),
		Team:    fb.team,
		Pos:     robotPosition,
		Dest:    fb.targetPosition,
		Dribble: true,
	}
	return &act
}

func (m *MoveWithBallToPosition) Achieved(gi *info.GameInfo) bool {
	ballPosition, err := gi.State.GetBall().GetPosition()
	if err != nil {
		Logger.Errorf("Position retrieval failed - Ball: %v\n", err)
		return false
	}
	robotPosition, err := gi.State.GetRobot(m.id, m.team).GetPosition()
	if err != nil {
		Logger.Errorf("Position retrieval failed - Robot: %v\n", err)
		return false
	}

	distanceLeft := ballPosition.Distance(m.targetPosition)
	const distance_threshold = 100
	const angle_threshold = 0.1
	distance_achieved := distanceLeft <= distance_threshold

	angle_diff := robotPosition.AngleDistance(m.targetPosition)
	angle_achieved := angle_diff <= angle_threshold
	return distance_achieved && angle_achieved
}

func (m *MoveWithBallToPosition) GetID() info.ID {
	return m.id
}
