package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/LiU-SeeGoals/controller/internal/gamestate"
)

func RunPythonHelloWorld() {
	resp, err := http.Get("http://controller-python:5000/")
	if err != nil {
		// println(err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			}
			bodyString := string(bodyBytes)
			println(bodyString)
		}
	}
}

func SendGameState(gs *gamestate.GameState) {
	output := gs.ToJson()
	reader := bytes.NewReader(output)
	resp, err := http.Post("http://controller-python:5000/update_game_state/", "application/json", reader)

	if err != nil {
		fmt.Println(err)
	} else {
		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		} else {
			print(string(bodyBytes))
		}
	}
}
