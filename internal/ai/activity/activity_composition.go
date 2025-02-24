package ai

import (
	"math"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type ActivityComposition struct {
	team info.Team
	id   info.ID
}

// Here we have funciton that are common across multiple activities,
// such as calculating the distance between two points.
// or movement that is legal and not blocked by other players.

func (fb *MoveWithBallToPosition) MoveWithBallToPosition(pos info.Position, gi *info.GameInfo) action.Action {
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
	act.Dest = pos

	// Reduce distance when its possible to estimte invisible ball
	if distance > 1500 {
		act.Dest = ballPos
	} else {
		act.Dest = pos
	}
	return &act
}

func (fb *MoveToBall) MoveToBall(gi *info.GameInfo) action.Action {
	myTeam := gi.State.GetTeam(fb.team)
	robot := myTeam[fb.id]
	if !robot.IsActive() {
		return nil
	}
	act := action.MoveTo{}
	act.Id = int(robot.GetID())
	act.Team = fb.team
	act.Pos = robot.GetPosition()
	act.Dest = gi.State.GetBall().GetPosition()
	act.Dribble = false
	return &act
}

func calculateDistance(p1, p2 info.Position) float32 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	dz := p1.Z - p2.Z
	return float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
}
