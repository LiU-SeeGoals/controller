package ai

import (
	"fmt"
	"math"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/info"
	. "github.com/LiU-SeeGoals/controller/internal/logger"
)

type MoveWithBallToPosition struct {
	GenericComposition
	target_position info.Position
	possession      bool
}

func (m *MoveWithBallToPosition) String() string {
	return fmt.Sprintf("Robot %d, MoveWithBallToPosition(%v)", m.id, m.target_position)
}

func NewMoveWithBallToPosition(team info.Team, id info.ID, dest info.Position) *MoveWithBallToPosition {
	return &MoveWithBallToPosition{
		GenericComposition: GenericComposition{
			team: team,
			id:   id,
		},
		target_position: dest,
	}
}

// TODO: Assumes that ones the robot has been in possession of the ball,
// it will keep it. It also doesnt support being initially in possession,
// for it to do that we need to estimate which robot has possession of the ball.
// But this might be enough to start work with other activities which depends
// on ball manipulation.
func (fb *MoveWithBallToPosition) GetAction(gi *info.GameInfo) action.Action {

	myTeam := gi.State.GetTeam(fb.team)
	robot := myTeam[fb.id]
	if !robot.IsActive() {
		return nil
	}

	robotPos, err1 := robot.GetPosition()

	if err1 != nil {
		Logger.Errorf("Position retrieval failed - Robot: %v\n", err1)
		return NewStop(fb.id).GetAction(gi)
	}

	ballPos, _, err := gi.State.GetBall().GetPositionTime()
	if err != nil {
		Logger.Errorf("Position retrieval failed - Ball: %v\n", err)
		return NewStop(fb.id).GetAction(gi)
	}
	// timeSinceUpdate := time.Now().UnixMilli() - updated

	// Check if ball position is upToDate
	// it might be old if the ball is
	// covered by the robot
	// upToDate := true
	// Logger.Debug("MoveWithBallToPosition: Time since update: %d", timeSinceUpdate)
	// if timeSinceUpdate > 0 {
	// 	upToDate = false
	// }

	distance := robotPos.Distance(ballPos)

	// If ball is far away and we are sure of
	// its position, move to it
	if distance > 100 && !fb.possession {
		Logger.Debug("MoveWithBallToPosition: Ball is far away")

		move := NewMoveToBall(fb.team, fb.id)
		return move.GetAction(gi)

	// If ball is close and we are not facing it
	// move to face it
	} else if !robot.FacingPosition(ballPos, 0.1) && !fb.possession {
		Logger.Debug("MoveWithBallToPosition: Ball is close and not facing it")

		dest := robotPos
		dest.Angle = robotPos.AngleToPosition(ballPos)
		move := NewMoveToPosition(fb.team, fb.id, dest)
		return move.GetAction(gi)

	}

	Logger.Debug("MoveWithBallToPosition: Moving with ball to position")
	// Passed all checks, is in possession
	// move with ball to target position
	fb.possession = true
	act := action.MoveTo{
		Id:      int(fb.id),
		Team:    fb.team,
		Pos:     robotPos,
		Dest:    fb.target_position,
		Dribble: true,
	}
	return &act
}

func (m *MoveWithBallToPosition) Achieved(gi *info.GameInfo) bool {
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
	const angle_threshold = 0.1
	distance_achieved := distance_left <= distance_threshold
	angle_diff := math.Abs(float64(curr_pos.Angle - target_position.Angle))
	angle_achieved := angle_diff <= angle_threshold
	return distance_achieved && angle_achieved
}

func (m *MoveWithBallToPosition) SetTargetPosition(dest info.Position) {
	m.target_position = dest
}
  
func (m *MoveWithBallToPosition) GetID() info.ID {
	return m.id
}

