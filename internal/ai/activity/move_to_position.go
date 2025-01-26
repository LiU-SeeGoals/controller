package ai

import (
	"math"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type MoveToPosition struct {
	team            info.Team
	id              info.ID
	target_position info.Position
}

func NewMoveToPosition(team info.Team, id info.ID, dest info.Position) *MoveToPosition {
	return &MoveToPosition{
		team:            team,
		id:              id,
		target_position: dest,
	}
}

func (m *MoveToPosition) GetAction(gi *info.GameInfo) action.Action {
	act := action.MoveTo{}
	act.Id = int(m.id)
	act.Team = m.team
	act.Pos = gi.State.GetTeam(m.team)[m.id].GetPosition()
	act.Dest = m.target_position
	act.Dribble = false
	return &act
}

func (m *MoveToPosition) Achieved(gi *info.GameInfo) bool {
	curr_pos := gi.State.GetTeam(m.team)[m.id].GetPosition()
	distance_left := calculateDistance(curr_pos, m.target_position)
	const distance_threshold = 10
	const angle_threshold = 0.1
	distance_achieved := distance_left <= distance_threshold
	angle_diff := math.Abs(float64(curr_pos.Angle - m.target_position.Angle))
	angle_achieved := angle_diff <= angle_threshold
	return distance_achieved && angle_achieved
}

func calculateDistance(p1, p2 info.Position) float32 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	dz := p1.Z - p2.Z
	return float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
}
