package ai

import (
	"fmt"
	"time"

	"github.com/LiU-SeeGoals/controller/internal/height_map"
	"github.com/LiU-SeeGoals/controller/internal/search"
	"github.com/LiU-SeeGoals/controller/internal/state"
)

type SlowBrainGO struct {
	team                 state.Team
	gameAnalysis         *state.GameAnalysis
	gameSearch           *search.FindBestScore
	incomingGameState    <-chan state.GameState
	outgoingPlan         chan<- state.GamePlan
	myAccumulatedFunc    height_map.HeightMap
	otherAccumulatedFunc height_map.HeightMap
}

func NewSlowBrainGO() SlowBrainGO {
	return SlowBrainGO{}
}

func (sb *SlowBrainGO) Init(incoming <-chan state.GameState, outgoing chan<- state.GamePlan, team state.Team) {
	posFunc := func(r *state.RobotAnalysis) state.Position {
		return r.GetPosition()
	}

	destFunc := func(r *state.RobotAnalysis) state.Position {
		return r.GetDestination()
	}

	myTimeAdvantage := height_map.NewTimeAdvantage(destFunc)
	otherTimeAdvantage := height_map.NewTimeAdvantage(posFunc)

	myAccumulatedFunc := func(x float32, y float32, robots state.RobotAnalysisTeam) float32 {
		scoreFuncs := []height_map.HeightMap{
			myTimeAdvantage.CalculateTimeAdvantage,
		}
		accumulated := float32(0)
		for _, scoreFunc := range scoreFuncs {
			accumulated += scoreFunc(x, y, robots)
		}
		return accumulated
	}

	otherAccumulatedFunc := func(x float32, y float32, robots state.RobotAnalysisTeam) float32 {
		scoreFuncs := []height_map.HeightMap{
			otherTimeAdvantage.CalculateTimeAdvantage,
		}
		accumulated := float32(0)
		for _, scoreFunc := range scoreFuncs {
			accumulated += scoreFunc(x, y, robots)
		}
		return accumulated
	}
	gameAnalysis := state.NewGameAnalysis(9000, 6000, 100, team)
	gameSearch := search.NewFindBestScore(team, myAccumulatedFunc, 0.1, 100, 9)
	sb.team = team
	sb.incomingGameState = incoming
	sb.outgoingPlan = outgoing
	sb.myAccumulatedFunc = myAccumulatedFunc
	sb.otherAccumulatedFunc = otherAccumulatedFunc
	sb.gameAnalysis = gameAnalysis
	sb.gameSearch = gameSearch

	go sb.Run()
}

func (sb *SlowBrainGO) Run() {
	var gameState state.GameState
	for {
		gameState = <-sb.incomingGameState

		time.Sleep(1 * time.Second) // TODO: Remove this
		// Wait for the game to start
		if !gameState.IsValid() {
			fmt.Println("SlowBrainGO: Invalid game state")
			time.Sleep(10 * time.Millisecond)
			continue
		}

		// Do some thinking
		plan := sb.GetPlan(&gameState)

		// Send the plan to the fast brain
		// fmt.Println(plan.ToDTO())
		sb.outgoingPlan <- plan
		fmt.Println("SlowBrainGO: Sent plan")
	}
}

func (sb *SlowBrainGO) GetPlan(gameState *state.GameState) state.GamePlan {
	sb.gameAnalysis.UpdateState(gameState)
	sb.gameAnalysis.UpdateMyZones(sb.myAccumulatedFunc)
	sb.gameAnalysis.UpdateOtherZones(sb.otherAccumulatedFunc)
	sb.gameSearch.FindBestScore(sb.myAccumulatedFunc, &sb.gameAnalysis.MyTeam, sb.gameAnalysis)
	gamePlan := state.GamePlan{}
	gamePlan.Team = sb.team
	for idx, _ := range sb.gameAnalysis.MyTeam.Robots {
		robot := &sb.gameAnalysis.MyTeam.Robots[idx]
		gamePlan.Instructions = append(gamePlan.Instructions, state.RobotMove{
			Id:       robot.GetID(),
			Position: robot.GetDestination(),
		})
	}
	gamePlan.Valid = true
	return gamePlan

}
