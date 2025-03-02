package ai

import (
	"fmt"
	"math"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type MoveToPosition struct {
	GenericComposition
	// MovementComposition
	target_position info.Position
}

func NewMoveToPosition(team info.Team, id info.ID, dest info.Position) *MoveToPosition {
	fmt.Println("NewMoveToPosition", id)
	return &MoveToPosition{
		GenericComposition: GenericComposition{
			team: team,
			id:   id,
		},
		target_position: dest,
	}
}

func (m *MoveToPosition) GetAction(gi *info.GameInfo) action.Action {
	act := action.MoveTo{}
	act.Id = int(m.id)
	act.Team = m.team
	// robot := gi.State.GetTeam(m.team)[m.id]
	act.Pos = gi.State.GetTeam(m.team)[m.id].GetPosition()
	// if m.team == info.Yellow {
	// 	robot := gi.State.GetTeam(m.team)[m.id]
	// 	target_position := m.avoidCollision(robot, m.target_position, gi.State)
	// 	act.Dest = target_position
	// } else {
	// 	act.Dest = m.target_position
	// }
	act.Dest = m.target_position

	act.Dribble = false
	return &act
}

func (m *MoveToPosition) Achieved(gi *info.GameInfo) bool {
	curr_pos := gi.State.GetTeam(m.team)[m.id].GetPosition()
	distance_left := CalculateDistance(curr_pos, m.target_position)
	const distance_threshold = 10
	const angle_threshold = 0.1
	distance_achieved := distance_left <= distance_threshold
	angle_diff := math.Abs(float64(curr_pos.Angle - m.target_position.Angle))
	angle_achieved := angle_diff <= angle_threshold
	return distance_achieved && angle_achieved
}
