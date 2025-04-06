package main

import (
	"fmt"
	"os"

	"github.com/LiU-SeeGoals/controller/internal/demos"
	"github.com/LiU-SeeGoals/controller/internal/visualisation"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "webdev" {
			fmt.Println("Starting webdev...")
			demos.RwSimScenario()
			return
		}
	}

	go demos.AoSimScenario()
	// Visualiser NEEDS to run in main thread
	vis := visualisation.GetVisualiser()
	vis.Run()
}
