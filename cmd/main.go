package main

import (
	"github.com/LiU-SeeGoals/controller/internal/demos"

	. "github.com/LiU-SeeGoals/controller/internal/logger"
)

func main() {

	defer Logger.Sync() // flushes buffer, if any
	demos.SimulatedAnnealing()

}
