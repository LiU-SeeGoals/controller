package height_map

import (
	"fmt"
	"math"

	"github.com/LiU-SeeGoals/controller/internal/state"
)

type HeightMap func(x float32, y float32, robots *state.RobotAnalysisTeam) float32

type HeightMapGauss struct {
	std float32
}

func (h HeightMapGauss) CalculateHeight(x float32, y float32, robots *state.RobotAnalysisTeam) float32 {
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

func (h HeightMapGauss) HeightMapDonut(x float32, y float32, robots *state.RobotAnalysisTeam) float32 {
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

func (ta *TimeAdvantage) CalculateTimeAdvantage(x float32, y float32, robots *state.RobotAnalysisTeam) float32 {
	time := float32(math.MaxFloat32)

	for _, robot := range robots {
		if !robot.IsActive() {
			continue
		}
		// Calculate the distance to the zone
		pos := ta.retrieveFunc(&robot)
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

// Function to find the best (x, y) based on height maps with float64 precision
func FindLowestHeight(id int, r float32, s int, heightMaps []HeightMap, gs *state.GameState) (float32, float32) {
	// Convert r to float64 for higher precision in the calculations
	r64 := float64(r)

	// Variables to track the best position and the minimum height (use float64 for precision)
	var bestX, bestY float64
	minHeight := math.MaxFloat64 // Start with a very high value

	// Get the current position of the yellow team robot (convert to float64)
	x := float64(gs.Yellow_team[id].GetPosition().At(0, 0))
	y := float64(gs.Yellow_team[id].GetPosition().At(1, 0))

	// Loop over `s` samples, evenly distributed around a circle
	for i := 0; i < s; i++ {
		// Calculate angle for this sample (in radians)
		angle := 2 * math.Pi * float64(i) / float64(s)

		// Compute new sample position (xSample, ySample) with float64 precision
		xSample := x + r64*math.Cos(angle)
		ySample := y + r64*math.Sin(angle)

		// Sum height map values at this sampled position
		totalHeight := float64(0)
		for _, heightMap := range heightMaps {
			// Cast back to float32 for the height map calculation (height maps expect float32)
			totalHeight += float64(heightMap.CalculateHeight(float32(xSample), float32(ySample), gs))
		}
		fmt.Println("totalHeight: ", xSample, ySample, totalHeight)

		// If this total height is lower than the current minimum, update best position
		if totalHeight < minHeight {
			minHeight = totalHeight
			bestX = xSample
			bestY = ySample
		}
	}

	// Divide bestX and bestY by 1000 to convert to millimeters
	bestX /= 1000
	bestY /= 1000

	// Return the best position (x, y) with the lowest height in millimeters
	// Convert back to float32 before returning
	return float32(bestX), float32(bestY)
}
