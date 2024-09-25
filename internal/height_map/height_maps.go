package height_map

import (
	"fmt"
	"math"

	"github.com/LiU-SeeGoals/controller/internal/gamestate"
)

type HeightMap interface {
	// x and y is the cordinate to avaluate the height
	// gamestate gives the context
	CalculateHeight(x float32, y float32, gs *gamestate.GameState) float32
}

type HeightMapEnemyGauss struct{}

func (h HeightMapEnemyGauss) CalculateHeight(x float32, y float32, gs *gamestate.GameState) float32 {
	// All enemy robots (blue team) create a positive Gaussian distribution
	gaussHeight := float32(0)
	stdDev := float64(1000) // Standard deviation is 1000

	for _, robot := range gs.Blue_team {
		if robot != nil && robot.GetPosition() != nil {
			// Get the robot's position
			robotX := float32(robot.GetPosition().At(0, 0))
			robotY := float32(robot.GetPosition().At(1, 0))

			// Gaussian falloff based on distance to the robot
			distanceSq := distanceSquared(x, y, robotX, robotY)

			// Gaussian with standard deviation 1000: adjust the denominator to 2 * stdDev^2
			gaussianValue := float32(math.Exp(-float64(distanceSq) / (2 * stdDev * stdDev)))

			gaussHeight += gaussianValue
		}
	}

	return gaussHeight
}

type HeightMapDonut struct{}

func (h HeightMapEnemyGauss) HeightMapDonut(x float32, y float32, gs *gamestate.GameState) float32 {
	// Around the robot closest to the ball in our team (yellow),
	// create a negative donut shaped distrobution at x distance
	return 0.5
}

type HeightMapAwayFromEdge struct{}

func (h HeightMapEnemyGauss) HeightMapAwayFromEdge(x float32, y float32, gs *gamestate.GameState) float32 {
	// It is often not advantagues to be close to the corneds of the playing field,
	// this creates incentive to not be close to corners.
	// The playing field is (-3,-4.5) to (3,4.5) in dimentions
	return 0.5
}

//------------------------------------------------------//
//   		Helper functions							//
//------------------------------------------------------//

func distanceSquared(x1, y1, x2, y2 float32) float32 {
	return (x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)
}

// Function to find the best (x, y) based on height maps with float64 precision
func FindLowestHeight(id int, r float32, s int, heightMaps []HeightMap, gs *gamestate.GameState) (float32, float32) {
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
