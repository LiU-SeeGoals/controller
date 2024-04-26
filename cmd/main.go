package main

import (
	"fmt"
	"os"

	"github.com/LiU-SeeGoals/controller/internal/ai"
	"github.com/LiU-SeeGoals/controller/internal/config"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"github.com/LiU-SeeGoals/controller/internal/world_predictor"
	"github.com/joho/godotenv"
)

func main() {
	conf, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	// go webserver.Once.Do(webserver.StartWebServer)

	var grsim_addr = "grsim:10300" //= conf.GRSIM_ADDR + ":" + conf.GRSIM_COMMAND_LISTEN_PORT
	// grsim_addr1 := "grsim:10301"
	// var grsim_addr2 := "grsim:10302"
	var vision = conf.SSL_VISION_MULTICAST_ADDR + ":" + conf.SSL_VISION_MAIN_PORT

	fmt.Println(grsim_addr)
	fmt.Println(vision)

	gs := gamestate.NewGameState(grsim_addr, vision)
	worldPredictor := world_predictor.NewWorldPredictor(config.GetSSLClientAddress(), gs)
	ai := ai.NewAi(gs, "grsim:10302")
	for {

		worldPredictor.Update()
		ai.Update()
		// fmt.Println(gs)
	}

}

// package main

// import (
// 	"fmt"
// 	"os"

// 	"github.com/LiU-SeeGoals/controller/internal/gamestate"
// 	"github.com/LiU-SeeGoals/controller/internal/webserver"
// 	"github.com/joho/godotenv"
// )

// func main() {
// 	conf, err := LoadConfig()
// 	if err != nil {
// 		panic(err)
// 	}

// 	go webserver.Once.Do(webserver.StartWebServer)

// 	var grsim_addr = conf.GRSIM_ADDR + ":" + conf.GRSIM_COMMAND_LISTEN_PORT
// 	var vision = conf.SSL_VISION_MULTICAST_ADDR + ":" + conf.SSL_VISION_MAIN_PORT

// 	fmt.Println(grsim_addr)
// 	fmt.Println(vision)

// 	gs := gamestate.NewGameState(grsim_addr, vision)

// 	//testAction := createTestActionMove(gs)
// 	// testAction2 := createTestActionInit(gs)
// 	for {
// 		//gs.AddAction(testAction)
// 		// gs.AddAction(testAction2)
// 		gs.Update()

// 		//fmt.Println(gs)
// 	}
// }

// func createTestActionMove(gs *gamestate.GameState) action.Action {
// 	id := 0
// 	act := &action.Move{}
// 	act.Pos = gs.GetRobot(id, gamestate.Yellow).GetPosition()
// 	act.Dest = mat.NewVecDense(3, nil)
// 	act.Id = id
// 	act.Dest.SetVec(0, 0)
// 	act.Dest.SetVec(1, 0)
// 	act.Dest.SetVec(2, 0)
// 	return act
// }

type Config struct {
	SSL_VISION_MAIN_PORT      string
	GRSIM_ADDR                string
	GRSIM_COMMAND_LISTEN_PORT string
	SSL_VISION_MULTICAST_ADDR string
}

func LoadConfig() (*Config, error) {
	// Load .env file. Adjust the path according to your .env file location.
	// If the .env file is in the same directory as the main package, you can use godotenv.Load() without arguments.
	if err := godotenv.Load("../.env"); err != nil {
		return nil, err
	}

	// Create config structure and populate it
	config := &Config{
		SSL_VISION_MAIN_PORT:      os.Getenv("SSL_VISION_MAIN_PORT"),
		GRSIM_ADDR:                os.Getenv("GRSIM_ADDR"),
		GRSIM_COMMAND_LISTEN_PORT: os.Getenv("GRSIM_COMMAND_LISTEN_PORT"),
		SSL_VISION_MULTICAST_ADDR: os.Getenv("SSL_VISION_MULTICAST_ADDR"),
	}

	return config, nil
}
