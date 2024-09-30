package ai

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

type PreCalculator struct {
	analysis *GameAnalysis
}

type Zone struct {
	timeAdvantage            float64 // time advantage of the zone
	anticipatedTimeAdvantage float64 // anticipated time advantage of the zone
}

// Struct to hold the analysis of the gamestate
type GameAnalysis struct {
	team        state.Team
	zones       [][]Zone // 2D array of zones
	fieldLength float32
	fieldWidth  float32
	zoneSize    float32
}

// Constructor for the PreCalculator
func NewPreCalculator(fieldLength, fieldWidth, zoneSize float32, team state.Team) *PreCalculator {
	pc := &PreCalculator{
		analysis: newAnalysis(fieldLength, fieldWidth, zoneSize, team),
	}
	return pc
}

// GameAnalysis constructor
func newAnalysis(fieldLength, fieldWidth, zoneSize float32, team state.Team) *GameAnalysis {
	analysis := GameAnalysis{}
	higth := int(fieldLength / zoneSize)
	width := int(fieldWidth / zoneSize)
	analysis.team = team
	analysis.fieldLength = fieldLength
	analysis.fieldWidth = fieldWidth
	analysis.zoneSize = zoneSize
	zones := make([][]Zone, higth)

	// Initialize the zones
	for i := 0; i < higth; i++ {
		zones[i] = make([]Zone, width)
	}
	analysis.zones = zones
	return &analysis
}

func (an *GameAnalysis) calculateTime(robots [state.TEAM_SIZE]*state.Robot, i, j int, fun func(*state.Robot) *mat.VecDense) float64 {
	time := math.Inf(1)
	// midel of the playfield in 0,0 so the zone need to be adjusted to the correct position
	posX := float32(i)*an.zoneSize - an.fieldLength/2 + an.zoneSize/2
	posY := float32(j)*an.zoneSize - an.fieldWidth/2 + an.zoneSize/2
	for _, robot := range robots {
		// Calculate the distance to the zone
		robotPos := fun(robot)
		rX := robotPos.AtVec(0)
		rY := robotPos.AtVec(1)
		zoneToRobot := mat.NewVecDense(2, []float64{float64(posX) - rX, float64(posY) - rY})
		distance := mat.Norm(zoneToRobot, 2)
		// Calculate the time to reach the zone
		curr_time := distance / robot.GetSpeed()
		if time > curr_time {
			time = curr_time
		}
	}
	return time
}

func (an *GameAnalysis) calculateTimeAdvantage(gamestateObj *state.GameState, i, j int, fun func(*state.Robot) *mat.VecDense) float64 {
	timeYellow := an.calculateTime(gamestateObj.GetYellowRobots(), i, j, fun)
	timeBlue := an.calculateTime(gamestateObj.GetBlueRobots(), i, j, fun)

	if an.team == state.Yellow {
		return timeBlue - timeYellow
	} else {
		return timeYellow - timeBlue
	}

}

func (an *GameAnalysis) updateTimeAdvantage(gamestateObj *state.GameState) {
	pos_func := func(r *state.Robot) *mat.VecDense {
		return r.GetPosition()
	}

	for i := 0; i < len(an.zones); i++ {
		for j := 0; j < len(an.zones[i]); j++ {
			// Calculate the time advantage of the zone
			an.zones[i][j].timeAdvantage = an.calculateTimeAdvantage(gamestateObj, i, j, pos_func)
		}
	}
}

func (an *GameAnalysis) updateAntisipetedTimeAdvantage(gamestateObj *state.GameState) {
	anticipate_func := func(r *state.Robot) *mat.VecDense {
		return r.GetAnticipatedPosition()
	}
	for i := 0; i < len(an.zones); i++ {
		for j := 0; j < len(an.zones[i]); j++ {
			// Calculate the time advantage of the zone
			an.zones[i][j].anticipatedTimeAdvantage = an.calculateTimeAdvantage(gamestateObj, i, j, anticipate_func)
		}
	}
}

func (pc *PreCalculator) Analyse(gamestateObj *state.GameState) *GameAnalysis {
	pc.analysis.updateTimeAdvantage(gamestateObj)

	return pc.analysis
}
