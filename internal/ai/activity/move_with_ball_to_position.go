package ai

import (
	"math"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type MoveWithBallToPosition struct {
	team            info.Team
	id              info.ID
	target_position info.Position
}

func NewMoveWithBallToPosition(team info.Team, id info.ID, dest info.Position) *MoveToPosition {
	return &MoveToPosition{
		team:            team,
		id:              id,
		target_position: dest,
	}
}

func (fb *MoveWithBallToPosition) GetAction(inst *info.Instruction, gs *info.GameState) action.Action {
	myTeam := gs.GetTeam(fb.team)
	robot := myTeam[inst.Id]
	if !robot.IsActive() {
		return nil
	}

	robotPos := robot.GetPosition()
	ballPos, _ := gs.GetBall().GetPositionTime()
	dx := float64(robotPos.X - ballPos.X)
	dy := float64(robotPos.Y - ballPos.Y)
	distance := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))

	act := action.MoveTo{}
	act.Id = int(robot.GetID())
	act.Team = fb.team
	act.Pos = robot.GetPosition()
	act.Dribble = true
	act.Dest = inst.Position

	// Reduce distance when its possible to estimte invisible ball
	if distance > 1500 {
		act.Dest = ballPos
	} else {
		act.Dest = inst.Position
	}
	return &act
}
