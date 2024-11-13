package search

import (
	"math"

	"github.com/LiU-SeeGoals/controller/internal/height_map"
	"github.com/LiU-SeeGoals/controller/internal/state"
)

// Function to find the best (x, y) based on height maps with float64 precision
func FindLowestHeight(x, y, r float32, s int, hightFunc height_map.HeightMap, robots state.RobotAnalysisTeam) (float32, float32) {
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
