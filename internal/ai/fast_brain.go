package ai

import (
	"fmt"
	"math"
	"time"

	. "github.com/LiU-SeeGoals/controller/internal/logger"
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
			Logger.Warn("FastBrainGO: Invalid game state")
			// fmt.Println("FastBrainGO: Invalid game state")
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

func (fb *FastBrainGO) Goalie(inst *info.Instruction, gs *info.GameState) action.Action {
	// todo: add collision avoidance
	myTeam := gs.GetTeam(fb.team)
	robot := myTeam[inst.Id]
	shooter := myTeam[1] 
	
	robotPos := robot.GetPosition()
	shooterPos := shooter.GetPosition()
	if !robot.IsActive() {
		return nil
	}
	act := action.MoveTo{}
	act.Id = int(robot.GetID())
	act.Team = fb.team
	act.Pos = robotPos

	ballPos, _ := gs.GetBall().GetPositionTime()
	
	//dright := math.Abs(float64(shooterPos.Y - 500))

	
	//drightgoalie := math.Abs(float64(shooterPos.Y - robotPos.Y))

	//angleMiddle := math.Atan(float64(math.abs(shooterPos.Y)/2000))

	//dleft := math.Abs(float64(shooterPos.Y + 500))
			//dleftgoalie := math.Abs(float64(shooterPos.Y + i))

		//angleGoalie := math.Atan(float64(dleftgoalie/2000))
		//anglePose := math.Atan(float64(dleft/2000))

	if(ballPos.X <= 4000){ 	
		if(shooterPos.Y <= 0){
			if(shooterPos.Y <= -500){
				act.Dest.Y = -350
			}else if(shooterPos.Y <= -350){
				act.Dest.Y = -250
			}else if(shooterPos.Y <= -250){
				act.Dest.Y = -150
			}else{
				act.Dest.Y = shooterPos.Y
			}
			
			

		}else{
			if(shooterPos.Y >= 500){
				act.Dest.Y = 350
			}else if(shooterPos.Y >= 350){
				act.Dest.Y = 250
			}else if(shooterPos.Y >= 250){
				act.Dest.Y = 150
			}else{
				act.Dest.Y = shooterPos.Y
			}
		}
	
		act.Dest.X = 4000
	}else{
		
		act.Dest.Y =  ballPos.Y
		act.Dest.X = ballPos.X+25
	}	

	act.Dribble = false
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
		Logger.Error("FastBrainGO: KickToPlayer not implemented")
	} else if inst.Type == info.KickToGoal {
		Logger.Error("FastBrainGO: KickToGoal not implemented")
	} else if inst.Type == info.KickToPosition{
		Logger.Error("FastBrainGO:KickToPosition not implemented")
	} else if inst.Type == info.ReceiveBallFromPlayer{
		Logger.Error("FastBrainGO:ReceiveBallFromPlayer not implemented")
	} else if inst.Type == info.ReceiveBallAtPosition{
		Logger.Error("FastBrainGO:ReceiveBallAtPosition not implemented")
	} else if inst.Type == info.BlockEnemyPlayerFromPosition {
		Logger.Error("FastBrainGO: BlockEnemyPlayerFromPosition not implemented")
	} else if inst.Type == info.BlockEnemyPlayerFromBall {
		Logger.Error("FastBrainGO: BlockEnemyPlayerFromBall not implemented")
	} else if inst.Type == info.BlockEnemyPlayerFromGoal {
		Logger.Error("FastBrainGO: BlockEnemyPlayerFromGoal not implemented")
	} else if inst.Type == info.BlockEnemyPlayerFromPlayer {
		Logger.Error("FastBrainGO: BlockEnemyPlayerFromPlayer not implemented")
	} else if inst.Type == info.Goalie {
		return fb.Goalie(inst, gs)
	} else {
		Logger.Error("FastBrainGO: not implemented")
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
