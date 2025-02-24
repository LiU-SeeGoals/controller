package ai

import (
	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type MoveWithBallToPosition struct {
	ActivityComposition
	team            info.Team
	id              info.ID
	target_position info.Position
}

func NewMoveWithBallToPosition(team info.Team, id info.ID, dest info.Position) *MoveToPosition {
	return &MoveToPosition{
		ActivityComposition: ActivityComposition{
			team: team,
			id:   id,
		},
		target_position: dest,
	}
}

func (fb *MoveWithBallToPosition) GetAction(gi *info.GameInfo) action.Action {

	pos := fb.target_position
	return fb.MoveWithBallToPosition(pos, gi)
}

func (fb *MoveWithBallToPosition) Achieved(*info.GameInfo) bool {
	// Need to be implemented
	return false
}
