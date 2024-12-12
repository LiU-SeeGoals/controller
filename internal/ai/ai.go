package ai

import (
	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/helper"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type SlowBrain interface {
	Init(incoming <-chan info.GameState, outgoing chan<- info.GamePlan, team info.Team)
}

type FastBrain interface {
	Init(incoming <-chan info.GameState, incomingPlan <-chan info.GamePlan, outgoing chan<- []action.Action, team info.Team)
}

type Ai struct {
	team              info.Team
	slow_brain        SlowBrain
	fast_brain        FastBrain
	gameStateSenderSB chan<- info.GameState
	gameStateSenderFB chan<- info.GameState
	actionReceiver    chan []action.Action
}

// Constructor for the ai, initializes the client
// and the different components used in the decision pipeline
func NewAi(team info.Team, slowBrain SlowBrain, fastBrain FastBrain) *Ai {
	gameStateSenderSB, gameStateReceiverSB := helper.NB_KeepLatestChan[info.GameState]()
	gameStateSenderFB, gameStateReceiverFB := helper.NB_KeepLatestChan[info.GameState]()
	gamePlanSender, gamePlanReceiver := helper.NB_KeepLatestChan[info.GamePlan]()
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
func (ai *Ai) GetActions(gi *info.GameInfo) []action.Action {

	// Send the game state copy to the slow brain
	ai.gameStateSenderSB <- *gi.State

	// Send the game state to the fast brain
	ai.gameStateSenderFB <- *gi.State

	// Get the actions from the fast brain, this will block until the fast brain has decided on actions
	actions := <-ai.actionReceiver

	return actions
}
