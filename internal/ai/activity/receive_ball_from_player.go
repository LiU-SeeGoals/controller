package ai

import (
	"math"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type ReceiveBallFromPlayer struct {
	ActivityComposition
	team     info.Team
	id       info.ID
	other_id info.ID
}

func NewReceiveBallFromPlayer(team info.Team, id info.ID, other_id info.ID) *ReceiveBallFromPlayer {
	return &ReceiveBallFromPlayer{
		ActivityComposition: ActivityComposition{
			team: team,
			id:   id,
		},
		other_id: other_id,
	}
}

func (fb *ReceiveBallFromPlayer) GetAction(gi *info.GameInfo) action.Action {

	myTeam := gi.State.GetTeam(fb.team)
	robotReceiver := myTeam[fb.id]
	if !robotReceiver.IsActive() {
		return nil
	}
	robotKicker := myTeam[fb.other_id]
	receiverPos := robotReceiver.GetPosition()
	kickerPos := robotKicker.GetPosition()

	ballPos, _ := gi.State.GetBall().GetPositionTime()
	dxBall := float64(receiverPos.X - ballPos.X)
	dyBall := float64(receiverPos.Y - ballPos.Y)
	distanceBall := math.Sqrt(math.Pow(dxBall, 2) + math.Pow(dyBall, 2))

	dx := float64(kickerPos.X - receiverPos.X)
	dy := float64(kickerPos.Y - receiverPos.Y)
	distance := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))

	if distanceBall < (distance / 3) {
		move := &MoveToBall{
			team: fb.team,
			id:   fb.id,
		}
		moveAction := move.MoveToBall(gi)
		moveAction.(*action.MoveTo).Dribble = true
		return moveAction
	}

	targetAngle := math.Atan2(math.Abs(dyBall), math.Abs(dxBall))
	if dx > 0 {
		targetAngle = math.Pi - targetAngle
	}
	if dy > 0 {
		targetAngle = -targetAngle
	}

	//because opposit angle
	if targetAngle > 0 {
		targetAngle -= math.Pi
	} else {
		targetAngle += math.Pi
	}

	//Rotate towards the kicker
	move := &MoveWithBallToPosition{
		team: fb.team,
		id:   fb.id,
	}
	pos := info.Position{X: receiverPos.X, Y: receiverPos.Y, Z: receiverPos.Z, Angle: float32(targetAngle)}
	return move.MoveWithBallToPosition(pos, gi)

	//Also needs to fix so that it moves out of the way if there is an obsticle

}

func (fb *ReceiveBallFromPlayer) Achieved(gi *info.GameInfo) bool {
	ballPos, _ := gi.State.GetBall().GetPositionTime()
	receiverPos := gi.State.GetTeam(fb.team)[fb.id].GetPosition()
	distance := calculateDistance(ballPos, receiverPos)
	const distance_threshold = 10
	ballRecived := distance <= distance_threshold
	return ballRecived
}
