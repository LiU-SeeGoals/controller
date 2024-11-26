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

func NewFastBrainGO() *FastBrainGO {
	return &FastBrainGO{}
}

func (fb *FastBrainGO) Init(incomingGameState <-chan state.GameState, incomingGamePlan <-chan state.GamePlan, outgoingActions chan<- []action.Action, team state.Team) {

	fb.incomingGameState = incomingGameState
	fb.incomingGamePlan = incomingGamePlan
	fb.outgoingActions = outgoingActions
	fb.team = team
	//

	go fb.Run()
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
		// time.Sleep(1 * time.Second) // TODO: Remove this

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
		// fmt.Println("FastBrainGO: Sent actions")

	}
}

func (fb *FastBrainGO) moveToPosition(inst *state.Instruction, gs *state.GameState) action.Action {
	// todo: add collision avoidance
	myTeam := gs.GetTeam(fb.team)
	robot := myTeam[inst.Id]
	if !robot.IsActive() {
		return nil
	}
	act := action.MoveTo{}
	act.Id = int(robot.GetID())
	act.Team = fb.team
	act.Pos = robot.GetPosition()
	if fb.team == state.Yellow {
		act.Dest = avoidObstacles(robot, inst.Position, *gs)
	} else {
		act.Dest = inst.Position
	}
	act.Dribble = false
	return &act
}

// TODO: can we make this nicer?
func (fb *FastBrainGO) instructionToAction(inst *state.Instruction, gs *state.GameState) action.Action {
	if inst.Type == state.MoveToPosition {
		return fb.moveToPosition(inst, gs)
	} else if inst.Type == state.MoveToBall {
		fmt.Println("FastBrainGO: MoveToBall not implemented")
	} else if inst.Type == state.MoveWithBallToPosition {
		fmt.Println("FastBrainGO: MoveWithBallToPosition not implemented")
	} else if inst.Type == state.KickToPlayer {
		fmt.Println("FastBrainGO: KickToPlayer not implemented")
	} else if inst.Type == state.KickToGoal {
		fmt.Println("FastBrainGO: KickToGoal not implemented")
	} else if inst.Type == state.KickToPosition {
		fmt.Println("FastBrainGO: KickToPosition not implemented")
	} else if inst.Type == state.ReceiveBallFromPlayer {
		fmt.Println("FastBrainGO: ReceiveBallFromPlayer not implemented")
	} else if inst.Type == state.ReceiveBallAtPosition {
		fmt.Println("FastBrainGO: ReceiveBallAtPosition not implemented")
	} else if inst.Type == state.BlockEnemyPlayerFromPosition {
		fmt.Println("FastBrainGO: BlockEnemyPlayerFromPosition not implemented")
	} else if inst.Type == state.BlockEnemyPlayerFromBall {
		fmt.Println("FastBrainGO: BlockEnemyPlayerFromBall not implemented")
	} else if inst.Type == state.BlockEnemyPlayerFromGoal {
		fmt.Println("FastBrainGO: BlockEnemyPlayerFromGoal not implemented")
	} else if inst.Type == state.BlockEnemyPlayerFromPlayer {
		fmt.Println("FastBrainGO: BlockEnemyPlayerFromPlayer not implemented")
	} else {
		fmt.Println("FastBrainGO: not implemented")
	}
	return nil
}

func (fb *FastBrainGO) GetActions(gs *state.GameState, gamePlan *state.GamePlan) []action.Action {

	var actionList []action.Action

	if fb.team != gamePlan.Team {
		panic("FastBrainGO: Team mismatch")
	}

	Instructions := gamePlan.Instructions

	for _, inst := range Instructions {
		action := fb.instructionToAction(inst, gs)
		if action != nil {
			actionList = append(actionList, action)
		}
	}

	return actionList
}
