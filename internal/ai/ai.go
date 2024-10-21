package ai

import (
	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/helper"
	"github.com/LiU-SeeGoals/controller/internal/state"
)

type SlowBrain interface {
	Init(incoming <-chan state.GameState, outgoing chan<- state.GamePlan, team state.Team)
}

type FastBrain interface {
	Init(incoming <-chan state.GameState, incomingPlan <-chan state.GamePlan, outgoing chan<- []action.Action, team state.Team)
}

type Ai struct {
	team              state.Team
	slow_brain        SlowBrain
	fast_brain        FastBrain
	gameStateSenderSB chan<- state.GameState
	gameStateSenderFB chan<- state.GameState
	actionReceiver    chan []action.Action
}

// Constructor for the ai, initializes the client
// and the different components used in the decision pipeline
func NewAi(team state.Team, slowBrain SlowBrain, fastBrain FastBrain) *Ai {
	gameStateSenderSB, gameStateReceiverSB := helper.NB_KeepLatestChan[state.GameState]()
	gameStateSenderFB, gameStateReceiverFB := helper.NB_KeepLatestChan[state.GameState]()
	gamePlanSender, gamePlanReceiver := helper.NB_KeepLatestChan[state.GamePlan]()
	actionReceiver := make(chan []action.Action)
	slowBrain.Init(gameStateReceiverSB, gamePlanSender, team)
	fastBrain.Init(gameStateReceiverFB, gamePlanReceiver, actionReceiver, team)
	ai := &Ai{
		team:              team,
		slow_brain:        slowBrain,
		fast_brain:        fastBrain,
		gameStateSenderSB: gameStateSenderSB,
		gameStateSenderFB: gameStateSenderFB,
		actionReceiver:    actionReceiver,
	}
	return ai
}

// Decides on new actions for the robots
func (ai *Ai) GetActions(gamestate *state.GameState) []action.Action {

	// Send the game state copy to the slow brain
	ai.gameStateSenderSB <- *gamestate

	// Send the game state to the fast brain
	ai.gameStateSenderFB <- *gamestate

	// Get the actions from the fast brain, this will block until the fast brain has decided on actions
	actions := <-ai.actionReceiver

	return actions
}
