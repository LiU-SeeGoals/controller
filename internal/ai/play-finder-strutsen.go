package ai

import (
	"math/rand"

	"github.com/LiU-SeeGoals/controller/internal/gamestate"
)

func NewPlayFinder() *StrategyFinder {
	pf := &StrategyFinder{}
	return pf
}

type StrategyFinder struct {
}

func fieldControleScore(field [][]Zone, fun func(*Zone) float64) float64 {
	posetive := 0.0
	negative := 0.0
	for i := 0; i < len(field); i++ {
		for j := 0; j < len(field[i]); j++ {
			if fun(&field[i][j]) > 0 {
				posetive++
			} else {
				negative++
			}
		}
	}

	return posetive / (posetive + negative)

}

func (pf *StrategyFinder) FindStrategy(gamestateObj *gamestate.GameState, gameAnalysis *GameAnalysis) float64 {
	currTimeFunc := func(zone *Zone) float64 {
		return zone.timeAdvantage
	}

	anticipatedTimeFunc := func(zone *Zone) float64 {
		return zone.anticipatedTimeAdvantage
	}
	anticipatedScore := 0.0

	currentScore := fieldControleScore(gameAnalysis.zones, currTimeFunc)
	gamestateObj.ResetAnticipetedPositions()

	myTeam := gamestateObj.GetTeam(gameAnalysis.team)
	for _, value := range rand.Perm(len(myTeam)) {
		// try this max 10 times
		for i := 0; i < 10; i++ {
			robot := myTeam[value]
			pos := robot.GetPosition()
			x := pos.AtVec(0)
			y := pos.AtVec(1)
			w := pos.AtVec(2)
			// TODO: anticipe position in a circle around the robot not a square
			dX := (rand.Float64() - 0.5) * 2 * robot.GetSpeed()
			dY := (rand.Float64() - 0.5) * 2 * robot.GetSpeed()
			robot.SetAnticipatedPosition(x+dX, y+dY, w)

			gameAnalysis.updateAntisipetedTimeAdvantage(gamestateObj)
			anticipatedScore = fieldControleScore(gameAnalysis.zones, anticipatedTimeFunc)

			if anticipatedScore > currentScore {
				break
			}
			robot.ResetAnticipatePosition()
		}
		break
	}
	return anticipatedScore
}
