package ai

import (
	"fmt"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/info"
	. "github.com/LiU-SeeGoals/controller/internal/logger"
)

type KickTheBall struct {
	GenericComposition
	targetPosition info.Position
	retrievingBall bool
}

func (k *KickTheBall) String() string {
	return fmt.Sprintf("(Robot %d, KickTheBall(%d))", k.id, k.targetPosition)
}

func NewKickTheBall(team info.Team, id info.ID, targetPosition info.Position) *KickTheBall {
	return &KickTheBall{
		GenericComposition: GenericComposition{
			team: team,
			id:   id,
		},
		targetPosition: targetPosition,
	}
}

func (kp *KickTheBall) GetAction(gi *info.GameInfo) action.Action {
	robot := gi.State.GetRobot(kp.id, kp.team)

	if !kp.retrievingBall { // Check if it lost the ball
		kp.retrievingBall = gi.State.LostBall(robot)
	}

	move := NewMoveToBall(kp.team, kp.id)
	if kp.retrievingBall && move.Achieved(gi) { // We have achivied in retrieving the ball
		Logger.Debug("MoveWithBallToPosition: Ball retrieved")
		kp.retrievingBall = false

	} else if kp.retrievingBall { // We are still working on getting the ball
		Logger.Debug("MoveWithBallToPosition: Retrieving ball")
		return move.GetAction(gi)

	}

	// Ensure we are close enough to the ball before attempting a kick
	robotPos, errR := robot.GetPosition()
	ballPos, errB := gi.State.GetBall().GetEstimatedPosition()
	if errR != nil || errB != nil {
		Logger.Errorf("Position retrieval failed - Robot/Ball: %v %v\n", errR, errB)
		return NewStop(kp.id).GetAction(gi)
	}

	distanceToBall := robotPos.Distance(ballPos)
	if distanceToBall > 90 { // WARN: Magic number, must be close to control the ball for a proper kick
		return move.GetAction(gi)
	}

	// Aim at target before kicking
	if !robot.Facing(kp.targetPosition, 0.1) {
		Logger.Debug("KickTheBall: Aiming at target before kick")
		angleToTarget := robotPos.AngleToPosition(kp.targetPosition)
		robotPos.Angle = angleToTarget
		return NewMoveWithBallToPosition(kp.team, kp.id, robotPos).GetAction(gi)
	}
	action := action.Kick{
		Id:        int(kp.id),
		KickSpeed: 4,
	}

	return &action
}

func (k *KickTheBall) Achieved(gi *info.GameInfo) bool {
	robot := gi.State.GetRobot(k.id, k.team)
	fmt.Println("Kicked away the ball: ", gi.State.LostBall(robot))
	return gi.State.LostBall(robot)
}

func (k *KickTheBall) GetID() info.ID {
	return k.id
}
