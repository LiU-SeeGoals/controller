package ai

import (
	"math"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type MoveToBall struct {
	GenericComposition
}

func NewMoveToBall(team info.Team, id info.ID) *MoveToPosition {
	return &MoveToPosition{
		GenericComposition: GenericComposition{
			team: team,
			id:   id,
		},
	}
}

func (m *MoveToBall) GetAction(gi *info.GameInfo) action.Action {
	target_position := gi.State.GetBall().GetPosition()
	act := action.MoveTo{}
	act.Id = int(m.id)
	act.Team = m.team
	act.Pos = gi.State.GetTeam(m.team)[m.id].GetPosition()
	act.Dest = target_position
	act.Dribble = false
	return &act
}

func (m *MoveToBall) Achieved(gi *info.GameInfo) bool {
	target_position := gi.State.GetBall().GetPosition()
	curr_pos := gi.State.GetTeam(m.team)[m.id].GetPosition()
	distance_left := CalculateDistance(curr_pos, target_position)
	const distance_threshold = 10
	const angle_threshold = 0.1
	distance_achieved := distance_left <= distance_threshold
	angle_diff := math.Abs(float64(curr_pos.Angle - target_position.Angle))
	angle_achieved := angle_diff <= angle_threshold
	return distance_achieved && angle_achieved
}
