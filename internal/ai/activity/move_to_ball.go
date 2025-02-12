package ai

import (
	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type MoveToBall struct {
	team            info.Team
	id              info.ID
	target_position info.Position
}

func NewMoveToBall(team info.Team, id info.ID, dest info.Position) *MoveToPosition {
	return &MoveToPosition{
		team:            team,
		id:              id,
		target_position: dest,
	}
}

func (fb *MoveToBall) GetAction(inst *info.Instruction, gs *info.GameState) action.Action {
	myTeam := gs.GetTeam(fb.team)
	robot := myTeam[inst.Id]
	if !robot.IsActive() {
		return nil
	}
	act := action.MoveTo{}
	act.Id = int(robot.GetID())
	act.Team = fb.team
	act.Pos = robot.GetPosition()
	act.Dest = gs.GetBall().GetPosition()
	act.Dribble = false
	return &act
}

func (m *MoveToBall) Achieved(gs *info.GameState) bool {
	// Need to be implemented
	return false
}
