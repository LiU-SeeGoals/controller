package ai

import (
	"time"
	"github.com/LiU-SeeGoals/controller/internal/state"
)

type SlowBrainGO struct {
	team 			state.Team
	incomingGameState <-chan state.GameState
	outgoingPlan      chan<- state.GamePlan
}

func NewSlowBrainGO(incoming <-chan state.GameState, outgoing chan<- state.GamePlan, team state.Team) *SlowBrainGO {
	sb := &SlowBrainGO{
		team: 			  team,
		incomingGameState: incoming,
		outgoingPlan:      outgoing,
	}
	return sb
}

func (sb *SlowBrainGO) Run() {
	var gameState state.GameState
	for {
		gameState = <-sb.incomingGameState

		// Wait for the game to start
		if gameState.Valid == false {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// Do some thinking
		plan := sb.GetPlan(&gameState)

		// Send the plan to the AI
		sb.outgoingPlan <- plan
	}
}

func (sb *SlowBrainGO) GetPlan(gameState *state.GameState) state.GamePlan {


