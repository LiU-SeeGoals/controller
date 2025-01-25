package ai

import (
	"fmt"
	"math"
	"time"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type FastBrainGO struct {
	team              info.Team
	incomingGameState <-chan info.GameState
	incomingGamePlan  <-chan info.GamePlan
	outgoingActions   chan<- []action.Action
}

func NewFastBrainGO() *FastBrainGO {
	return &FastBrainGO{}
}

func (fb *FastBrainGO) Init(incomingGameState <-chan info.GameState, incomingGamePlan <-chan info.GamePlan, outgoingActions chan<- []action.Action, team info.Team) {

	fb.incomingGameState = incomingGameState
	fb.incomingGamePlan = incomingGamePlan
	fb.outgoingActions = outgoingActions
	fb.team = team
	//

	go fb.Run()
}

func (fb *FastBrainGO) Run() {
	gameState := info.GameState{}
	gamePlan := info.GamePlan{}

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

func (fb *FastBrainGO) moveToPosition(inst *info.Instruction, gs *info.GameState) action.Action {
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
	act.Dest = inst.Position
	// if fb.team == info.Yellow {
	// 	act.Dest = avoidObstacles(robot, inst.Position, *gs)
	// } else {
	// 	act.Dest = inst.Position
	// }
	act.Dribble = false
	return &act
}

func (fb *FastBrainGO) moveToBall(inst *info.Instruction, gs *info.GameState) action.Action {
	myTeam := gs.GetTeam(fb.team)
	robot := myTeam[inst.Id]
	if !robot.IsActive() {
		return nil
	}
	act := action.MoveTo{}
	act.Id = int(robot.GetID())
	act.Team = fb.team
	act.Pos = robot.GetPosition()
	act.Dest = gs.GetBall().GetPosition()
	act.Dribble = false
	return &act
}

func (fb *FastBrainGO) moveWithBallToPosition(inst *info.Instruction, gs *info.GameState) action.Action {
	myTeam := gs.GetTeam(fb.team)
	robot := myTeam[inst.Id]
	if !robot.IsActive() {
		return nil
	}

	robotPos := robot.GetPosition()
	ballPos, _ := gs.GetBall().GetPositionTime()
	dx := float64(robotPos.X - ballPos.X)
	dy := float64(robotPos.Y - ballPos.Y)
	distance := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))

	act := action.MoveTo{}
	act.Id = int(robot.GetID())
	act.Team = fb.team
	act.Pos = robot.GetPosition()
	act.Dribble = true
	act.Dest = inst.Position

	// Reduce distance when its possible to estimte invisible ball
	if distance > 1500 {
		act.Dest = ballPos
	} else {
		act.Dest = inst.Position
	}
	return &act
}

// TODO: can we make this nicer?
func (fb *FastBrainGO) instructionToAction(inst *info.Instruction, gs *info.GameState) action.Action {
	if inst.Type == info.MoveToPosition {
		return fb.moveToPosition(inst, gs)
	} else if inst.Type == info.MoveToBall {
		return fb.moveToBall(inst, gs)
	} else if inst.Type == info.MoveWithBallToPosition {
		return fb.moveWithBallToPosition(inst, gs)
		//fmt.Println("FastBrainGO: MoveWithBallToPosition not implemented")
	} else if inst.Type == info.KickToPlayer {
		fmt.Println("FastBrainGO: KickToPlayer not implemented")
	} else if inst.Type == info.KickToGoal {
		fmt.Println("FastBrainGO: KickToGoal not implemented")
	} else if inst.Type == info.KickToPosition {
		fmt.Println("FastBrainGO: KickToPosition not implemented")
	} else if inst.Type == info.ReceiveBallFromPlayer {
		fmt.Println("FastBrainGO: ReceiveBallFromPlayer not implemented")
	} else if inst.Type == info.ReceiveBallAtPosition {
		fmt.Println("FastBrainGO: ReceiveBallAtPosition not implemented")
	} else if inst.Type == info.BlockEnemyPlayerFromPosition {
		fmt.Println("FastBrainGO: BlockEnemyPlayerFromPosition not implemented")
	} else if inst.Type == info.BlockEnemyPlayerFromBall {
		fmt.Println("FastBrainGO: BlockEnemyPlayerFromBall not implemented")
	} else if inst.Type == info.BlockEnemyPlayerFromGoal {
		fmt.Println("FastBrainGO: BlockEnemyPlayerFromGoal not implemented")
	} else if inst.Type == info.BlockEnemyPlayerFromPlayer {
		fmt.Println("FastBrainGO: BlockEnemyPlayerFromPlayer not implemented")
	} else {
		fmt.Println("FastBrainGO: not implemented")
	}
	return nil
}

func (fb *FastBrainGO) GetActions(gs *info.GameState, gamePlan *info.GamePlan) []action.Action {

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
