package ai

import (
	"fmt"
	"math"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type MoveWithBallToPosition struct {
	GenericComposition
	target_position info.Position
}

func (m *MoveWithBallToPosition) String() string {
	return fmt.Sprintf("MoveWithBallToPosition(%d, %d, %v)", m.team, m.id, m.target_position)
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

func (fb *MoveWithBallToPosition) GetAction(gi *info.GameInfo) action.Action {

	myTeam := gi.State.GetTeam(fb.team)
	robot := myTeam[fb.id]
	if !robot.IsActive() {
		return nil
	}

	robotPos := robot.GetPosition()
	ballPos, _ := gi.State.GetBall().GetPositionTime()
	dx := float64(robotPos.X - ballPos.X)
	dy := float64(robotPos.Y - ballPos.Y)
	distance := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))

	act := action.MoveTo{}
	act.Id = int(robot.GetID())
	act.Team = fb.team
	act.Pos = robot.GetPosition()
	act.Dribble = true
	act.Dest = fb.target_position

	// Reduce distance when its possible to estimte invisible ball
	if distance > 1500 {
		act.Dest = ballPos
	} else {
		act.Dest = fb.target_position
	}
	return &act
}

func (m *MoveWithBallToPosition) Achieved(gi *info.GameInfo) bool {
	target_position := gi.State.GetBall().GetPosition()
	curr_pos := gi.State.GetTeam(m.team)[m.id].GetPosition()
	distance_left := curr_pos.Distance(target_position)
	const distance_threshold = 100
	const angle_threshold = 0.1
	distance_achieved := distance_left <= distance_threshold
	angle_diff := math.Abs(float64(curr_pos.Angle - target_position.Angle))
	angle_achieved := angle_diff <= angle_threshold
	return distance_achieved && angle_achieved
}
