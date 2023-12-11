package main

import (
	"fmt"

	"github.com/LiU-SeeGoals/controller/internal/config"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
)

func main() {

	gs := gamestate.NewGameState(config.GetGrSimAddress(), config.GetSSLClientAddress())
	for {
		gs.TestActions()
		gs.Update()
		fmt.Println(gs)
	}
}
