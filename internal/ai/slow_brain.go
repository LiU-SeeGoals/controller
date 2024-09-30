package ai

import (
	"github.com/LiU-SeeGoals/controller/internal/state"
)

type SlowBrainGO struct {
	incomingGameState chan<- state.GameState
	outgoingPlan      <-chan state.GamePlan
}

func NewSlowBrainGO(incoming chan<- state.GameState, outgoing <-chan state.GamePlan) *SlowBrainGO {
	sb := &SlowBrainGO{
		incomingGameState: incoming,
		outgoingPlan:      outgoing,
	}
	return sb
}

func (sb *SlowBrainGO) Run() {
	for {
		gameState := <-sb.incomingGameState

		// Do some thinking
		plan := sb.GetPlan(&gameState)

		// Send the plan to the AI
		sb.outgoingPlan <- plan
	}
}
