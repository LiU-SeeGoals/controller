package search

import (
	"math"

	"github.com/LiU-SeeGoals/controller/internal/height_map"
	"github.com/LiU-SeeGoals/controller/internal/state"
)

type FindRngBestScore struct {
	team         state.Team
	gameAnalysis *state.GameAnalysis
}

func NewFindRngBestScore(team state.Team, gameAnalysis *state.GameAnalysis) *FindRngBestScore {
	return &FindRngBestScore{
		team:         team,
		gameAnalysis: gameAnalysis,
	}
}

func timeAdvantage(posX, posY float32, robot *state.RobotAnalysis) float32 {
	pos := robot.GetPosition()
	robotX := pos.X
	robotY := pos.Y

	dist := float32(math.Sqrt(math.Pow(float64(posX-robotX), 2) + math.Pow(float64(posY-robotY), 2)))

	time := dist / robot.GetMaxMoveSpeed()

	return time
}

func scoreTimeAdvantageZones(an state.GameAnalysis) float32 {
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

func MinScore(robot *state.RobotAnalysis,
	an *state.GameAnalysis,
	fieldInfo *state.FieldInfo,
	zoneScoreFunc func(float32, float32, *state.RobotAnalysis) float32,
	analysisScoreFunc func(state.GameAnalysis) float32,
) float32 {
	myTeam := an.MyTeam

	for i := 0; i < len(myTeam.Zones); i++ {
		posX := float32(i)*myTeam.ZoneSize - fieldInfo.Length/2 + myTeam.ZoneSize/2
		for j := 0; j < len(myTeam.Zones[i]); j++ {
			// middle of the play field in 0,0 so the zone need to be adjusted to the correct position
			posY := float32(j)*myTeam.ZoneSize - fieldInfo.Width/2 + myTeam.ZoneSize/2
			zoneScore := float32(math.MaxFloat32)
			for _, robot := range myTeam.Robots {
				if !robot.IsActive() {
					continue
				}
				score := zoneScoreFunc(posX, posY, &robot)
				if score < zoneScore {
					zoneScore = score
				}
			}
			myTeam.Zones[i][j].Score = zoneScore
		}
	}

	return analysisScoreFunc(*an)
}

func (frbs *FindRngBestScore) FindRngBestScore(x, y, r float32, s int, hightFunc height_map.HeightMap, robots *state.RobotAnalysisTeam) (float32, float32) {
	// Variables to track the best position and the minimum height
	var bestX, bestY float32
	minHeight := float32(math.MaxFloat32) // Start with a very high value

	// Loop over `s` samples, evenly distributed around a circle
	for i := 0; i < s; i++ {
		// Calculate angle for this sample (in radians)
		angle := 2 * math.Pi * float64(i) / float64(s)

		// Compute new sample position (xSample, ySample) with float64 precision
		xSample := x + r*float32(math.Cos(angle))
		ySample := y + r*float32(math.Sin(angle))

		// Sum height map values at this sampled position
		totalHeight := hightFunc(xSample, ySample, robots)

		if totalHeight < minHeight {
			minHeight = totalHeight
			bestX = xSample
			bestY = ySample
		}
	}

	return bestX, bestY
}
