package ai

import (
	"fmt"
	"time"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/state"
	"gonum.org/v1/gonum/mat"
)

type FastBrainGO struct {
	team              state.Team
	incomingGameState <-chan state.GameState
	incomingGamePlan  <-chan state.GamePlan
	outgoingActions   chan<- []action.Action
}

func NewFastBrain[FB FastBrainGO](incomingGameState <-chan state.GameState, incomingGamePlan <-chan state.GamePlan, outgoingActions chan<- []action.Action, team state.Team) *FastBrainGO {

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

	for _, inst := range Instructions {

		act := action.MoveTo{}
		act.Id = int(inst.Id)

		robot := myTeam[inst.Id]

		pos := robot.GetPosition()

		act.Pos = mat.NewVecDense(3, []float64{float64(pos.X), float64(pos.Y), float64(pos.Angel)})
		dest := inst.Position
		act.Dest = mat.NewVecDense(3, []float64{float64(dest.X), float64(dest.Y), float64(dest.Angel)})

		act.Dribble = true // Assuming all moves require dribbling
		if act.Dest.AtVec(0) == act.Pos.AtVec(0) && act.Dest.AtVec(1) == act.Pos.AtVec(1) {
			continue
		}
		// fmt.Println("Robot", act.Id, "moving to", destX, destY, "from", act.Pos.AtVec(0), act.Pos.AtVec(1))
		actionList = append(actionList, &act)
	}

	return actionList
}
