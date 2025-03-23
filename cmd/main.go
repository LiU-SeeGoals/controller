package main

import (
	"fmt"
	"os"

	"github.com/LiU-SeeGoals/controller/internal/demos"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "webdev" {
			fmt.Println("Starting webdev...")
			demos.RwSimScenario()
			return
		}
	}

	demos.FwRealScenario()
}
