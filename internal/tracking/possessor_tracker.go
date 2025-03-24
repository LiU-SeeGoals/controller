package tracking

import (
	"math"

	"github.com/LiU-SeeGoals/controller/internal/info"
	. "github.com/LiU-SeeGoals/controller/internal/logger"
)

// Constants for the possession detector
const (
	maxDistance              = 95.0 // mm
	maxPrevPossessorDistance = 125  // mm
	facingThreshold          = 0.8  // radians
)

type possessorTracker struct {
	lastPossessor *info.Robot
}

func (tracker *possessorTracker) InPossession(robot *info.Robot, ball *info.Ball, lastPossessor *info.Robot) bool {

	ballPos, err := ball.GetPosition()
	if err != nil {
		return false
	}

	robotPos, err := robot.GetPosition()
	if err != nil {
		return false
	}

	// --- Distance ---
	// check distance for robots that are not the last possessor
	ballDistance := robotPos.Distance(ballPos)
	if ballDistance > maxDistance && lastPossessor != robot {
		return false
	}

	// // last possessor bias
	if ballDistance > maxPrevPossessorDistance && lastPossessor == robot {
		return false
	}

	// --- Facing ---
	if !robot.Facing(ballPos, facingThreshold) {
		return false
	}

	// // --- Velocity alignment ---
	ballVel := ball.GetVelocity()
	robotVel := robot.GetVelocity()

	// If ball is moving, robot should be moving roughly the same speed
	speedDiff := math.Abs(robotVel.Length() - ballVel.Length())
	if speedDiff > 0.5 {
		return false
	}

	// Robot and ball should be moving in the same direction
	direction := robotVel.Normalize().Dot(ballVel.Normalize())
	if direction < 0.5 {
		return false
	}

	return true
}

func (tracker *possessorTracker) GetPossessor(state *info.GameState) *info.Robot {

	ball := state.GetBall()

	var possessors []*info.Robot
	for _, team := range []info.Team{info.Yellow, info.Blue} {
		for _, robot := range state.GetTeam(team) {

			if tracker.InPossession(robot, ball, tracker.lastPossessor) {
				possessors = append(possessors, robot)
			}
		}
	}

	if len(possessors) == 0 {
		return nil

	} else if len(possessors) > 1 {
		Logger.Warnf("Multiple robots in possession, picking the first one: %v", possessors)
	}

	return possessors[0]
}
