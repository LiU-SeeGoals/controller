package ai

import (
	"fmt"
	"math"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/info"
	. "github.com/LiU-SeeGoals/controller/internal/logger"
)

type RamAtPosition struct {
	GenericComposition
	targetPosition info.Position
}

func (k *RamAtPosition) String() string {
	return fmt.Sprintf("(Robot %d, RamAtPosition(%d))", k.id, k.targetPosition)
}

func NewRamAtPosition(team info.Team, id info.ID, targetPosition info.Position) *RamAtPosition {
	return &RamAtPosition{
		GenericComposition: GenericComposition{
			team: team,
			id:   id,
		},
		targetPosition: targetPosition,
	}
}

func (kp *RamAtPosition) GetAction(gi *info.GameInfo) action.Action {
	robot := gi.State.GetRobot(kp.id, kp.team)
	robotPos, err := robot.GetPosition()
	if err != nil {
		Logger.Errorf("Position retrieval failed - Kicker: %v\n", err)
		return NewStop(kp.id).GetAction(gi)
	}

	ballPos, err := gi.State.Ball.GetPosition()
	if err != nil {
		Logger.Errorf("Position retrieval failed - Kicker: %v\n", err)
		return NewStop(kp.id).GetAction(gi)
	}

	// Stop at ball
	if robotPos.Distance(ballPos) < 50 {
		return NewStop(kp.id).GetAction(gi)
	}


	angleBallToStartPos := ballPos.AngleToPosition(kp.targetPosition) + math.Pi
	startPos := ballPos.OnRadius(500, angleBallToStartPos)

	move := NewMoveToPosition(kp.team, kp.id, startPos)
	move.AvoidBall(true)
	
	// In start position, RAM THE BALL
	if move.Achieved(gi) {
		angleBallToTargetPos := ballPos.AngleToPosition(kp.targetPosition)
		targetPos := ballPos.OnRadius(500, angleBallToTargetPos)
		return NewMoveToPosition(kp.team, kp.id, targetPos).GetAction(gi)
	}
	return move.GetAction(gi)
}

func (k *RamAtPosition) Achieved(gi *info.GameInfo) bool {
	robotPos, err := gi.State.GetRobot(k.id, k.team).GetPosition()
	if err != nil {
		Logger.Errorf("Position retrieval failed - Kicker: %v\n", err)
		return false
	}

	ballPos, err := gi.State.Ball.GetPosition()
	if err != nil {
		Logger.Errorf("Position retrieval failed - Kicker: %v\n", err)
		return false
	}
	return robotPos.Distance(ballPos) < 50
}

func (k *RamAtPosition) GetID() info.ID {
	return k.id
}
