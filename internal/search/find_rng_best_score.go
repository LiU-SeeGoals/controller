package search

import (
	"math"
	"math/rand"

	"github.com/LiU-SeeGoals/controller/internal/height_map"
	"github.com/LiU-SeeGoals/controller/internal/state"
)

type FindRngBestScore struct {
	team        state.Team
	scoringFunc func(x float32, y float32, robots state.RobotAnalysisTeam) float32
	alpha       float32
	radius      float32
	samples     int
}

func NewFindRngBestScore(
	team state.Team,
	scoringFunc func(x float32, y float32, robots state.RobotAnalysisTeam) float32,
	alpha float32,
	radius float32,
	samples int,
) *FindRngBestScore {
	return &FindRngBestScore{
		team:        team,
		scoringFunc: scoringFunc,
		alpha:       alpha,
		radius:      radius,
		samples:     samples,
	}
}

func scoreHighest(an *state.GameAnalysis) float32 {
	myScore := float32(0)
	otherScore := float32(0)

	myTeam := an.MyTeam
	otherTeam := an.OtherTeam

	for i := 0; i < len(myTeam.Zones); i++ {
		for j := 0; j < len(myTeam.Zones[i]); j++ {
			if myTeam.Zones[i][j].Score < otherTeam.Zones[i][j].Score {
				myScore += 1
			} else {
				otherScore += 1
			}
		}
	}

	return myScore / (myScore + otherScore)
}

func (frbs *FindRngBestScore) FindRngBestScore(hightFunc height_map.HeightMap, robotTeam state.TeamAnalysis, gameAnalysis *state.GameAnalysis) {
	// Variables to track the best position and the minimum height
	bestScore := scoreHighest(gameAnalysis)

	robots := robotTeam.Robots
	for _, robot := range robots {
		robotPos := robot.GetPosition()
		robot.SetDestination(robotPos)
	}
	for _, value := range rand.Perm(len(robots)) {
		robot := robots[value]
		pos := robot.GetPosition()
		x := pos.X
		y := pos.Y
		for i := 0; i < frbs.samples; i++ {
			// Calculate angle for this sample (in radians)
			angle := 2 * math.Pi * float64(i) / float64(frbs.samples)

			// Compute new sample position (xSample, ySample) with float64 precision
			xSample := x + frbs.radius*float32(math.Cos(angle))
			ySample := y + frbs.radius*float32(math.Sin(angle))

			dest := robot.GetDestination()
			dest.X = xSample
			dest.Y = ySample
			robot.SetDestination(dest)

			gameAnalysis.UpdateMyZones(frbs.scoringFunc)
			score := scoreHighest(gameAnalysis)
			if score > bestScore+frbs.alpha {
				bestScore = score
			} else {
				dest.X = x
				dest.Y = y
				robot.SetDestination(dest)
			}
		}
	}
}