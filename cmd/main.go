package main

import (
	// "github.com/LiU-SeeGoals/controller/internal/demos"
	"time"

	. "github.com/LiU-SeeGoals/controller/internal/logger"
)

func main() {
	for {
		time.Sleep(2000 * time.Millisecond)
		Logger.Info("This is an info message")
		Logger.Error("This is an error message")
		Logger.Debug("This is a debug message")
		Logger.Warn("This is a warning message")

	}

	//demos.PythonSlowBrain()
	//demos.Scenario()
	// demos.Goalie()
}
