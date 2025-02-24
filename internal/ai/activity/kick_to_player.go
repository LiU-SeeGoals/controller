package ai

import (
	"math"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type KickToPlayer struct {
	ActivityComposition
	team     info.Team
	id       info.ID
	other_id info.ID
}

func NewKickToPlayer(team info.Team, id info.ID, other_id info.ID) *KickToPlayer {
	return &KickToPlayer{
		ActivityComposition: ActivityComposition{
			team: team,
			id:   id,
		},
		other_id: other_id,
	}
}

func (fb *KickToPlayer) GetAction(gi *info.GameInfo) action.Action {
	myTeam := gi.State.GetTeam(fb.team)
	robotKicker := myTeam[fb.id]
	if !robotKicker.IsActive() {
		return nil
	}
	robotReciever := myTeam[fb.other_id]

	kickerPos := robotKicker.GetPosition()
	recieverPos := robotReciever.GetPosition()

	dx := float64(kickerPos.X - recieverPos.X)
	dy := float64(kickerPos.Y - recieverPos.Y)
	distance := math.Sqrt(dx*dx + dy*dy)

	targetAngle := math.Atan2(math.Abs(dy), math.Abs(dx))
	if dx > 0 {
		targetAngle = math.Pi - targetAngle
	}
	if dy > 0 {
		targetAngle = -targetAngle
	}

	ballPos, _ := gi.State.GetBall().GetPositionTime()
	dxBall := float64(kickerPos.X - ballPos.X)
	dyBall := float64(kickerPos.Y - ballPos.Y)
	distanceBall := math.Sqrt(math.Pow(dxBall, 2) + math.Pow(dyBall, 2))

	// Rotate to target
	if math.Abs(float64(kickerPos.Angle)-float64(targetAngle)) > 0.05 {
		move := &MoveWithBallToPosition{
			team: fb.team,
			id:   fb.id,
		}
		pos := info.Position{X: kickerPos.X, Y: kickerPos.Y, Z: kickerPos.Z, Angle: float32(targetAngle)}
		return move.MoveWithBallToPosition(pos, gi)
	}

	// kick
	if distanceBall > 90 {
		move := &MoveToBall{
			team: fb.team,
			id:   fb.id,
		}
		return move.MoveToBall(gi)
	} else {
		kickAct := &action.Kick{}
		kickAct.Id = int(robotKicker.GetID())

		// Compute the kick speed as a function of the distance to target
		normDistance := float64(distance) / 10816
		kickSpeed := 1 + int(4*normDistance)
		kickAct.KickSpeed = int(math.Min(math.Max(float64(kickSpeed), 1), 5))
		return kickAct
	}

	//Needs to add that is doesn't kick if there is an obsicle

}

func (k *KickToPlayer) Achieved(gi *info.GameInfo) bool {
	ballPos, _ := gi.State.GetBall().GetPositionTime()
	receiverPos := gi.State.GetTeam(k.team)[k.other_id].GetPosition()
	distance := calculateDistance(ballPos, receiverPos)
	const distance_threshold = 10
	ballRecived := distance <= distance_threshold
	return ballRecived
}
