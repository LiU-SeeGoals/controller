package ai

import (
	"time"

	"github.com/LiU-SeeGoals/controller/internal/height_map"
	"github.com/LiU-SeeGoals/controller/internal/state"
)

type SlowBrainGO struct {
	team                 state.Team
	gameAnalysis         *state.GameAnalysis
	incomingGameState    <-chan state.GameState
	outgoingPlan         chan<- state.GamePlan
	myAccumulatedFunc    height_map.HeightMap
	otherAccumulatedFunc height_map.HeightMap
}

func NewSlowBrainGO(incoming <-chan state.GameState, outgoing chan<- state.GamePlan, team state.Team) *SlowBrainGO {
	posFunc := func(r *state.RobotAnalysis) state.Position {
		return r.GetPosition()
	}

	destFunc := func(r *state.RobotAnalysis) state.Position {
		return r.GetDestination()
	}

	myTimeAdvantage := height_map.NewTimeAdvantage(destFunc)
	otherTimeAdvantage := height_map.NewTimeAdvantage(posFunc)

	myAccumulatedFunc := func(x float32, y float32, robots *state.RobotAnalysisTeam) float32 {
		scoreFuncs := []height_map.HeightMap{
			myTimeAdvantage.CalculateTimeAdvantage,
		}
		accumulated := float32(0)
		for _, scoreFunc := range scoreFuncs {
			accumulated += scoreFunc(x, y, robots)
		}
		return accumulated
	}

	otherAccumulatedFunc := func(x float32, y float32, robots *state.RobotAnalysisTeam) float32 {
		scoreFuncs := []height_map.HeightMap{
			otherTimeAdvantage.CalculateTimeAdvantage,
		}
		accumulated := float32(0)
		for _, scoreFunc := range scoreFuncs {
			accumulated += scoreFunc(x, y, robots)
		}
		return accumulated
	}

	sb := &SlowBrainGO{
		team:                 team,
		gameAnalysis:         state.NewGameAnalysis(9000, 6000, 100, team),
		incomingGameState:    incoming,
		outgoingPlan:         outgoing,
		myAccumulatedFunc:    myAccumulatedFunc,
		otherAccumulatedFunc: otherAccumulatedFunc,
	}
	go sb.Run()
	return sb
}

func (sb *SlowBrainGO) Run() {
	var gameState state.GameState
	for {
		gameState = <-sb.incomingGameState

		// Wait for the game to start
		if gameState.Valid == false {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// Do some thinking
		plan := sb.GetPlan(&gameState)

		// Send the plan to the fast brain
		sb.outgoingPlan <- plan
	}
}

func (sb *SlowBrainGO) GetPlan(gameState *state.GameState) state.GamePlan {
	sb.gameAnalysis.Update(gameState)

}
