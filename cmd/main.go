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

	// Visualiser (plotting) NEEDS to run in main thread
	// Choose backend, none gives no GUI fayne creates GUI
	// Use none if there are display issues, and you dont wanna debug
	vis := visualisation.NewVisualiser("none")
	vis.Run()
}
