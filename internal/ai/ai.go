package ai

import (
	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/helper"
	"github.com/LiU-SeeGoals/controller/internal/state"
)

type Ai struct {
	team              state.Team
	slow_brain        *SlowBrainGO
	fast_brain        *FastBrainGO
	gameStateSenderSB chan<- state.GameState
	gameStateSenderFB chan<- state.GameState
	actionReciver     chan []action.Action
}

// Constructor for the ai, initializes the client
// and the different components used in the decision pipeline
func NewAi(team state.Team) *Ai {
	gameStateSenderSB, gameStateRecivrerSB := helper.NB_KeepLatestChan[state.GameState]()
	gameStateSenderFB, gameStateRecivrerFB := helper.NB_KeepLatestChan[state.GameState]()
	gamePlanSender, gamePlanRecivrer := helper.NB_KeepLatestChan[state.GamePlan]()
	actionReciver := make(chan []action.Action)
	slowBrain := NewSlowBrainGO(gameStateRecivrerSB, gamePlanSender, team)
	fastBrain := NewFastBrainGO(gameStateRecivrerFB, gamePlanRecivrer, actionReciver, team)
	ai := &Ai{
		team:              team,
		slow_brain:        slowBrain,
		fast_brain:        fastBrain,
		gameStateSenderSB: gameStateSenderSB,
		gameStateSenderFB: gameStateSenderFB,
		actionReciver:     actionReciver,
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
	actions := <-ai.actionReciver

	return actions
}
