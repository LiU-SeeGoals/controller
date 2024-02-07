package main

import (
	"github.com/LiU-SeeGoals/controller/internal/world_predictor"
)

func main() {

	worldPredictor := world_predictor.NewWorldPredictor()
	worldPredictor.Update()
	//ai := ai.NewAi(gs, config.GetGrSimAddress())
	for {

		//worldPredictor.Update()
		//ai.Update()
	}

}
