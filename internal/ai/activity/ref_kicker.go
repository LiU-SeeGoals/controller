package ai

import (
	"fmt"
	"math"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type RefKicker struct {
	GenericComposition
	targetPos info.Position
}

func (m *RefKicker) String() string {
	return fmt.Sprintf("RefKicker(%d)", m.id)
}

func NewRefKicker(id info.ID, team info.Team, refState info.RefState) *RefKicker {
	return &RefKicker{
		GenericComposition: GenericComposition{
			id:   id,
			team: team,
		},
	}

}

func (m *RefKicker) GetAction(gi *info.GameInfo) action.Action {
	var act action.Action
	robotPos, _ := gi.State.GetRobot(m.id, m.team).GetPosition()
	ballPos, _ := gi.State.Ball.GetPosition()

	// Prepare kickoff, get in position
	if gi.Status.GetGameEvent().GetCurrentState() != info.STATE_KICKOFF_PREPARATION {
		targetPos := info.Position{X: -300, Y: 0, Z: 0, Angle: 0}
		if m.team == info.Blue && gi.Status.GetBlueTeamOnPositiveHalf() || m.team == info.Yellow && !gi.Status.GetBlueTeamOnPositiveHalf() {
			// We have the positive half
			targetPos = info.Position{X: 300, Y: 0, Z: 0, Angle: math.Pi}
		}
		m.targetPos = targetPos
		move := NewMoveToPosition(m.team, m.id, targetPos)
		move.AvoidBall(true)
		act = move.GetAction(gi)

	} else if gi.Status.GetGameEvent().GetCurrentState() != info.STATE_FREE_KICK {
		// "Kick" the ball

		targetPos := info.Position{X: 300, Y: 0, Z: 0, Angle: 0}
		if m.team == info.Blue && gi.Status.GetBlueTeamOnPositiveHalf() || m.team == info.Yellow && !gi.Status.GetBlueTeamOnPositiveHalf() {
			// We have the positive half
			targetPos = info.Position{X: -300, Y: 0, Z: 0, Angle: math.Pi}
		}
		act = NewRamAtPosition(m.team, m.id, targetPos).GetAction(gi)

	} 
	return act
}

func (m *RefKicker) Achieved(gi *info.GameInfo) bool {
	return gi.State.GetRobot(m.id, m.team).At(m.targetPos, 100)
}

func (m *RefKicker) GetID() info.ID {
	return m.id

}
