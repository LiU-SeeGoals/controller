package helper

import (
	"math/rand"

	"github.com/LiU-SeeGoals/controller/internal/info"
)

func Contains(slice []int, num int) bool {
	for _, item := range slice {
		if item == num {
			return true
		}
	}
	return false
}

// Generates a random position on the field
func Random_Position(gameinfo info.GameInfo) info.Position {
	// Position (0, 0) is at the center of the field.
	// Randomizes sign for each direction to allow for movement in all directions.
	// Therefor we cant move more than half the width and height of the field in any direction.
	var width = gameinfo.FieldSize().X * 0.7
	var height = gameinfo.FieldSize().Y * 0.7

	return info.Position{
		X:     rand.Float64() * (width / 2) * float64(RandSign()),
		Y:     rand.Float64() * (height / 2) * float64(RandSign()),
		Z:     0,
		Angle: 0,
	}
}

func RandSign() int {
	sign := 1
	if rand.Intn(2) == 0 {
		sign = -1
	}
	return sign
}
