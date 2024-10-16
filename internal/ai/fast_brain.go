package ai

import (
	"fmt"
	"time"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/state"
)

type FastBrainGO struct {
	team              state.Team
	incomingGameState <-chan state.GameState
	incomingGamePlan  <-chan state.GamePlan
	outgoingActions   chan<- []action.Action
}

func NewFastBrain(incomingGameState <-chan state.GameState, incomingGamePlan <-chan state.GamePlan, outgoingActions chan<- []action.Action, team state.Team) *FastBrainGO {

	fb := FastBrainGO{
		team:              team,
		incomingGameState: incomingGameState,
		incomingGamePlan:  incomingGamePlan,
		outgoingActions:   outgoingActions,
	}

	go fb.Run()
	return &fb
}

func (fb *FastBrainGO) Run() {
	gameState := state.GameState{}
	gamePlan := state.GamePlan{}
	for {
		// We will reive the game state more often than the game plan
		// so we wait for the gameState to update and work with the latest game plan

		gameState = <-fb.incomingGameState

		select {
		case gamePlan = <-fb.incomingGamePlan:
		default:

		}
		time.Sleep(1 * time.Second) // TODO: Remove this

		// Wait for the game to start
		if !gameState.Valid || !gamePlan.Valid {
			fmt.Println("FastBrainGO: Invalid game state")
			fb.outgoingActions <- []action.Action{}
			time.Sleep(10 * time.Millisecond)
			continue
		}

		// Do some thinking
		actions := fb.GetActions(&gameState, &gamePlan)

		// Send the actions to the AI
		fb.outgoingActions <- actions
		fmt.Println("FastBrainGO: Sent actions")

	}
}

func (fb *FastBrainGO) GetActions(gs *state.GameState, gamePlan *state.GamePlan) []action.Action {

	var actionList []action.Action

	myTeam := gs.GetTeam(gamePlan.Team)
	Instructions := gamePlan.Instructions

	for idx, _ := range Instructions {
		inst := &Instructions[idx]
		robot := &myTeam[inst.Id]

		if !robot.IsActive() {
			continue
		}
		act := action.MoveTo{}
		act.Id = int(inst.Id)

		pos := robot.GetPosition()

		act.Pos = pos
		act.Dest = inst.Position

		act.Dribble = true // Assuming all moves require dribbling
		fmt.Println("Robot", act.Id, "moving to", act.Dest.ToDTO())
		actionList = append(actionList, &act)
	}

	return actionList
}
