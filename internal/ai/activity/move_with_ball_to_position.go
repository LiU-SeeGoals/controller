package ai

import (
	"math"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type MoveWithBallToPosition struct {
	GenericComposition
	team            info.Team
	id              info.ID
	target_position info.Position
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

func (m *MoveWithBallToPosition) GetAction(gi *info.GameInfo) action.Action {
	robotPos := gi.State.GetTeam(m.team)[m.id].GetPosition()

	ballPos := gi.State.GetBall().GetPosition()
	dx := float64(robotPos.X - ballPos.X)
	dy := float64(robotPos.Y - ballPos.Y)
	distance := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))

	act := action.MoveTo{}
	act.Id = int(m.id)
	act.Team = m.team
	act.Pos = robotPos
	act.Dribble = true
	act.Dest = m.target_position

	// Reduce distance when its possible to estimte invisible ball
	if distance > 1500 {
		act.Dest = ballPos
	} else {
		act.Dest = m.target_position
	}
	return &act
}

func (m *MoveWithBallToPosition) Achieved(gi *info.GameInfo) bool {
	target_position := gi.State.GetBall().GetPosition()
	curr_pos := gi.State.GetTeam(m.team)[m.id].GetPosition()
	distance_left := curr_pos.Distance(target_position)
	const distance_threshold = 10
	const angle_threshold = 0.1
	distance_achieved := distance_left <= distance_threshold
	angle_diff := math.Abs(float64(curr_pos.Angle - target_position.Angle))
	angle_achieved := angle_diff <= angle_threshold
	return distance_achieved && angle_achieved
}
