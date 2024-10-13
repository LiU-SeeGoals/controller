package client

import (
	"bytes"
	"encoding/json"
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
				// println(err->)
			}
			bodyString := string(bodyBytes)
			println(bodyString)
		}
	}
}

func SendGameState(gs *gamestate.GameState) {
	output, err := json.Marshal(gs)
	if err != nil {
		// TODO: Do something
	}
	println(string(output))

	reader := bytes.NewReader(output)
	// _, err :=
	http.Post("http://controller-python:5000/update_game_state/", "application/json", reader)

	if err != nil {
		// TODO: Do something
	}

	// bodyBytes, err := io.ReadAll(resp.Body)

	// print(string(bodyBytes))
}
