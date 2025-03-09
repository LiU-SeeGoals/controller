package ai

import (
	"fmt"

	"github.com/LiU-SeeGoals/controller/internal/action"
	. "github.com/LiU-SeeGoals/controller/internal/logger"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type MoveToPosition struct {
	GenericComposition
	// MovementComposition
	target_position info.Position
}


func (m *MoveToPosition) String() string {
	return fmt.Sprintf("(Robot %d, MoveToPosition(%v))", m.id, m.target_position)
}

func NewMoveToPosition(team info.Team, id info.ID, dest info.Position) *MoveToPosition {
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
	robot := gi.State.GetTeam(m.team)[m.id]
	robotPos, err := robot.GetPosition()

	if err != nil {
		Logger.Errorf("Position retrieval failed - Robot: %v\n", err)
		return NewStop(m.id).GetAction(gi)
	}

	act.Pos = robotPos
	act.Dest = m.target_position


	act.Dribble = false
	return &act
}

func (m *MoveToPosition) Achieved(gi *info.GameInfo) bool {
	curr_pos, err := gi.State.GetTeam(m.team)[m.id].GetPosition()
	if err != nil {
		Logger.Errorf("Position retrieval failed - Robot: %v\n", err)
		return false
	}
	distance_left := curr_pos.Distance(m.target_position)
	const distance_threshold = 100
	const angle_threshold = 0.1
	distance_achieved := distance_left <= distance_threshold

	angle_diff := curr_pos.AngleDistance(m.target_position)
	angle_achieved := angle_diff <= angle_threshold
	return distance_achieved && angle_achieved
}

func (m *MoveToPosition) SetTargetPosition(dest info.Position) {
	m.target_position = dest
}

func (m *MoveToPosition) GetID() info.ID {
	return m.id
}

