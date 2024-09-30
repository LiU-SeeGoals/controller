package ai

import (
	"time"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/state"
)

type FastBrainGO struct {
	incomingGameState <-chan state.GameState
	incomingGamePlan  <-chan state.GamePlan
	outgoingActions   chan<- []action.Action
	gameState         state.GameState
	gamePlan          state.GamePlan
}

func NewFastBrainGO(incomingGameState <-chan state.GameState, incomingGamePlan <-chan state.GamePlan, outgoingActions chan<- []action.Action) *FastBrainGO {

	fb := &FastBrainGO{
		incomingGameState: incomingGameState,
		incomingGamePlan:  incomingGamePlan,
		outgoingActions:   outgoingActions,
		gameState:         state.GameState{},
		gamePlan:          state.GamePlan{},
	}
	go fb.Run()
	return fb
}

func (fb *FastBrainGO) Run() {
	for {
		// We will recive the game state more often than the game plan
		select {
		case fb.gameState = <-fb.incomingGameState:
		case fb.gamePlan = <-fb.incomingGamePlan:
		default:

		}

		// Wait for the game to start
		if !fb.gameState.Valid || !fb.gamePlan.Valid {
			time.Sleep(10 * time.Millisecond)
			continue
		}

		// Do some thinking
		actions := fb.GetActions(&gameState, &gamePlan)

		// Send the actions to the AI
		fb.outgoingActions <- actions
	}
}

func (fb *FastBrainGO) GetActions(gs *state.GameState, gameAnalysis *state.GamePlan) []action.Action {

	var actionList []action.Action

	myTeam := gs.GetTeam(gameAnalysis.team)

	for _, robot := range myTeam {

		act := action.MoveTo{}
		act.Pos = robot.GetPosition()
		act.Id = robot.GetID()

		anticipatePosition := robot.GetAnticipatedPosition()
		destX := anticipatePosition.AtVec(0)
		destY := anticipatePosition.AtVec(1)
		act.Dest = mat.NewVecDense(3, []float64{destX, destY, 0})

		act.Dribble = true // Assuming all moves require dribbling
		if destX == act.Pos.AtVec(0) && destY == act.Pos.AtVec(1) {
			continue
		}
		// fmt.Println("Robot", act.Id, "moving to", destX, destY, "from", act.Pos.AtVec(0), act.Pos.AtVec(1))
		actionList = append(actionList, &act)
	}

	return actionList
}
