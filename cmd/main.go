package main

import (
	"encoding/json"
	"os"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"github.com/LiU-SeeGoals/controller/internal/webserver"
	"gonum.org/v1/gonum/mat"
)

type Config struct {
	SSLClientAddress string `json:"sslClientAddress"`
	GrSimAddress     string `json:"grSimAddress"`
}

func main() {
	conf, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	go webserver.Once.Do(webserver.StartWebServer)

	gs := gamestate.NewGameState(conf.GrSimAddress, conf.SSLClientAddress)

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

// func createTestActionInit(gs *gamestate.GameState) action.Action {
// 	id := 1
// 	act := &action.Init{}
// 	act.Id = id
// 	return act
// }

func LoadConfig() (*Config, error) {
	file, err := os.Open("../config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
