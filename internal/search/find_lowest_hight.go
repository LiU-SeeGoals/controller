package search

import (
	"math"

	"github.com/LiU-SeeGoals/controller/internal/height_map"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

// Function to find the best (x, y) based on height maps with float64 precision
func FindLowestHeight(x, y, r float64, s int, hightFunc height_map.HeightMap, robots info.RobotAnalysisTeam) (float64, float64) {
	// Variables to track the best position and the minimum height
	var bestX, bestY float64
	minHeight := float64(math.MaxFloat32) // Start with a very high value

	// Loop over `s` samples, evenly distributed around a circle
	for i := 0; i < s; i++ {
		// Calculate angle for this sample (in radians)
		angle := 2 * math.Pi * float64(i) / float64(s)

		// Compute new sample position (xSample, ySample) with float64 precision
		xSample := x + r*float64(math.Cos(angle))
		ySample := y + r*float64(math.Sin(angle))

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
