// package main

// import (
// 	"fmt"

// 	"github.com/LiU-SeeGoals/controller/internal/ai"
// 	"github.com/LiU-SeeGoals/controller/internal/config"
// 	"github.com/LiU-SeeGoals/controller/internal/gamestate"
// 	"github.com/LiU-SeeGoals/controller/internal/world_predictor"
// )

// func main() {

// 	gs := gamestate.NewGameState()
// 	worldPredictor := world_predictor.NewWorldPredictor(config.GetSSLClientAddress(), gs)
// 	ai := ai.NewAi(gs, config.GetGrSimAddress())
// 	for {

// 		worldPredictor.Update()
// 		ai.Update()
// 		fmt.Println(gs)
// 	}

// }

package main

import (
	"os"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"github.com/LiU-SeeGoals/controller/internal/webserver"
	"github.com/joho/godotenv"
	"gonum.org/v1/gonum/mat"
)

func main() {
	conf, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	go webserver.Once.Do(webserver.StartWebServer)

	gs := gamestate.NewGameState(conf.GRSIM_ADDR, conf.SSL_VISION_MAIN_PORT)

	testAction := createTestActionMove(gs)
	// testAction2 := createTestActionInit(gs)
	for {
		gs.AddAction(testAction)
		// gs.AddAction(testAction2)
		gs.Update()

		//fmt.Println(gs)
	}
}

func createTestActionMove(gs *gamestate.GameState) action.Action {
	id := 0
	act := &action.Move{}
	act.Pos = gs.GetRobot(id, gamestate.Yellow).GetPosition()
	act.Dest = mat.NewVecDense(3, nil)
	act.Id = id
	act.Dest.SetVec(0, 0)
	act.Dest.SetVec(1, 0)
	act.Dest.SetVec(2, 0)
	return act
}

type Config struct {
	SSL_VISION_MAIN_PORT string
	GRSIM_ADDR           string
}

func LoadConfig() (*Config, error) {
	// Load .env file. Adjust the path according to your .env file location.
	// If the .env file is in the same directory as the main package, you can use godotenv.Load() without arguments.
	if err := godotenv.Load("../.env"); err != nil {
		return nil, err
	}

	// Create config structure and populate it
	config := &Config{
		SSL_VISION_MAIN_PORT: os.Getenv("SSL_VISION_MAIN_PORT"),
		GRSIM_ADDR:           os.Getenv("GRSIM_ADDR"),
	}

	return config, nil
}
