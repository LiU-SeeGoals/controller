package height_map

import (
	"math"

	"github.com/LiU-SeeGoals/controller/internal/state"
)

type HeightMap func(x float32, y float32, robots state.RobotAnalysisTeam) float32

type HeightMapGauss struct {
	std float32
}

func (h HeightMapGauss) CalculateHeight(x float32, y float32, robots state.RobotAnalysisTeam) float32 {
	// All enemy robots (blue team) create a positive Gaussian distribution
	gaussHeight := float32(0)

	for _, robot := range robots {
		if robot.IsActive() {
			pos := robot.GetPosition()
			robotX := pos.X
			robotY := pos.Y

			// Gaussian falloff based on distance to the robot
			distanceSq := distanceSquared(x, y, robotX, robotY)

			// Gaussian with standard deviation 1000: adjust the denominator to 2 * stdDev^2
			gaussianValue := float32(math.Exp(-float64(distanceSq) / float64(2*h.std*h.std)))

			gaussHeight += gaussianValue
		}
	}

	return gaussHeight
}

type HeightMapDonut struct{}

func (h HeightMapGauss) HeightMapDonut(x float32, y float32, robots state.RobotAnalysisTeam) float32 {
	// Around the robot closest to the ball in our team (yellow),
	// create a negative donut shaped distrobution at x distance
	return 0.5
}

type HeightMapAwayFromEdge struct{}

func (h HeightMapGauss) HeightMapAwayFromEdge(x float32, y float32, robots *state.RobotAnalysisTeam) float32 {
	// It is often not advantagues to be close to the corneds of the playing field,
	// this creates incentive to not be close to corners.
	// The playing field is (-3,-4.5) to (3,4.5) in dimentions
	return 0.5
}

type TimeAdvantage struct {
	retrieveFunc func(r *state.RobotAnalysis) state.Position
}

func NewTimeAdvantage(retrieveFunc func(r *state.RobotAnalysis) state.Position) *TimeAdvantage {
	return &TimeAdvantage{
		retrieveFunc: retrieveFunc,
	}
}

func (ta *TimeAdvantage) CalculateTimeAdvantage(x float32, y float32, robots state.RobotAnalysisTeam) float32 {
	time := float32(math.MaxFloat32)

	for _, robot := range robots {
		if !robot.IsActive() {
			continue
		}
		// Calculate the distance to the zone
		pos := ta.retrieveFunc(robot)
		pos.X = x - pos.X
		pos.Y = y - pos.Y
		distance := pos.Norm()
		// Calculate the time to reach the zone
		currTime := distance / robot.GetMaxMoveSpeed()
		if time > currTime {
			time = currTime
		}
	}
	return time
}

//------------------------------------------------------//
//   		Helper functions							//
//------------------------------------------------------//

func distanceSquared(x1, y1, x2, y2 float32) float32 {
	return (x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)
}
