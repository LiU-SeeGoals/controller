package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"time"
)

func main() {
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	fmt.Println("Use WASD to control the robot. Press 'ESC' to exit.")

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyEsc {
			break
		}

		switch char {
		case 'w':
			fmt.Println("Moving forward")
			// Add your code for 'w' key action (e.g., send a command to the robot)
		case 'a':
			fmt.Println("Moving left")
			// Add your code for 'a' key action
		case 's':
			fmt.Println("Moving backward")
			// Add your code for 's' key action
		case 'd':
			fmt.Println("Moving right")
			// Add your code for 'd' key action
		}

		// You may want to add a delay to avoid high CPU usage in a busy loop
		time.Sleep(100 * time.Millisecond)
	}
}
