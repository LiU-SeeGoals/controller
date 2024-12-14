package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/state"
	"github.com/eiannone/keyboard"
)

const (
	gsim = iota
	basestation
	remote_control
)

var (
	clientType   int = 0
	commands     map[rune]command
	speed        int  = 0
	robotStopped bool = false
)

type command struct {
	message string
	run     func() action.Action
}

func main() {
	clientHost := "127.0.0.1:20011"
	client := askForClient(clientHost)
	client.Init()
	initCommands(askForRobotId())
	listenKeyboard(client)
}

func askForClient(port string) client.Client {
	var userChoice string
	var clientHost string = port
	var clientBaseStation client.Client

	fmt.Println("Please enter the client type [g]sim (default), [b]ase station or [r]emote control: ")
	fmt.Scanln(&userChoice)

	if userChoice == "b" || userChoice == "r" {
		fmt.Println("Enter <ip>:<port> for the basestation (port defaults to 6001): ")
		fmt.Scanln(&clientHost)
		if !strings.Contains(clientHost, ":") {
			clientHost = clientHost + ":6001"
		}

		if userChoice == "b" {
			clientType = basestation
			fmt.Println("Creating base station client.")
		} else {
			clientType = remote_control
			fmt.Println("Creating base station client for remote control.")
		}
	} else {
		clientType = gsim
		fmt.Println("Creating gsim client.")
	}

	clientBaseStation = client.NewBaseStationClient(clientHost)
	return clientBaseStation
}

func askForRobotId() int {
	var robotId string

	fmt.Println("Please enter the robot ID (defaults to 0): ")
	fmt.Scanln(&robotId)
	id, err := strconv.Atoi(robotId)
	if err != nil {
		fmt.Println("Setting Robot ID to 0.")
		id = 0
	}
	fmt.Println("Robot ID is: ", id)
	return id
}

func initCommands(robotId int) {
	if clientType == gsim || clientType == basestation {
		commands = map[rune]command{
			'w': {
				message: "Moving forward",
				run: func() action.Action {
					return &action.Move{
						Id:        robotId,
						Direction: state.Position{X: 0.0, Y: 1.0},
					}
				},
			},
			'l': {
				message: "Stopping robot",
				run: func() action.Action {
					return &action.Stop{
						Id: robotId,
					}
				},
			},
			'a': {
				message: "Moving left",
				run: func() action.Action {
					return &action.Move{
						Id:        robotId,
						Direction: state.Position{X: -1.0, Y: 0.0},
					}
				},
			},
			's': {
				message: "Moving backward",
				run: func() action.Action {
					return &action.Move{
						Id:        robotId,
						Direction: state.Position{X: 0.0, Y: -1.0},
					}
				},
			},
			'd': {
				message: "Moving right",
				run: func() action.Action {
					return &action.Move{
						Id:        robotId,
						Direction: state.Position{X: 1.0, Y: 0.0},
					}
				},
			},
			'k': {
				message: "Kicking",
				run: func() action.Action {
					return &action.Kick{
						Id: robotId,
					}
				},
			},
		}
	} else { // remote control
		// In remote control mode, we've got some additional commands and expect
		// some special handling of other commands.
		commands = map[rune]command{
			'w': {
				message: "Moving forward",
				run: func() action.Action {
					robotStopped = false
					return &action.Move{
						Id:        robotId,
						Direction: state.Position{X: 0.0, Y: 1.0},
					}
				},
			},
			'a': {
				message: "Moving left",
				run: func() action.Action {
					robotStopped = false
					return &action.Move{
						Id:        robotId,
						Direction: state.Position{X: -1.0, Y: 0.0},
					}
				},
			},
			's': {
				message: "Moving backward",
				run: func() action.Action {
					robotStopped = false
					return &action.Move{
						Id:        robotId,
						Direction: state.Position{X: 0.0, Y: -1.0},
					}
				},
			},
			'd': {
				message: "Moving right",
				run: func() action.Action {
					robotStopped = false
					return &action.Move{
						Id:        robotId,
						Direction: state.Position{X: 1.0, Y: 0.0},
					}
				},
			},
			'l': {
				message: "Stopping robot",
				run: func() action.Action {
					robotStopped = true
					return &action.Stop{
						Id: robotId,
					}
				},
			},
			'k': {
				message: "Kicking",
				run: func() action.Action {
					robotStopped = false
					return &action.Kick{
						Id: robotId,
					}
				},
			},
			'q': {
				message: "Rotating left",
				run: func() action.Action {
					robotStopped = false
					return &action.Kick{
						Id: robotId,
					}
				},
			},
			'e': {
				message: "Rotating right",
				run: func() action.Action {
					robotStopped = false
					return &action.Kick{
						Id: robotId,
					}
				},
			},
			'r': {
				message: "Speed decreased",
				run: func() action.Action {
					speed -= 1
					fmt.Println(speed)
					return &action.Kick{
						Id:        robotId,
						KickSpeed: speed,
					}
				},
			},
			't': {
				message: "Speed increased",
				run: func() action.Action {
					speed += 1
					fmt.Println(speed)
					return &action.Kick{
						Id:        robotId,
						KickSpeed: speed,
					}
				},
			},
			'p': {
				message: "Sent ping",
				run: func() action.Action {
					return &action.Ping{
						Id: robotId,
					}
				},
			},
		}
	}
}

func sendCommand(char rune, client client.Client) {
	if cmd, exists := commands[char]; exists {
		fmt.Println(cmd.message)
		client.SendActions([]action.Action{cmd.run()})
	} else {
		fmt.Println("Bad command: ", char)
	}
}

func sendPing(client client.Client) {
	for {
		if !robotStopped {
			sendCommand('p', client)
		}
		time.Sleep(time.Second)
	}
}

func listenKeyboard(client client.Client) {
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	if clientType == gsim || clientType == basestation {
		fmt.Println("Use WASD to control the robot, <space> to stop all movement, K to kick.")
	} else {
		fmt.Println("Use WASD to control the robot, <space> to stop all movement, K to kick, O/P to decrease/increase speed.")
		fmt.Println("Pings are sent continually unless <space> is pressed.")
	}
	fmt.Println("Press <ESC> to exit.")

	// Send continous pings if we're remote controlling
	if clientType == remote_control {
		go sendPing(client)
	}

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyEsc {
			break
		} else if key == keyboard.KeySpace {
			// translate space to "stop" command
			char = 'l'
		}

		sendCommand(char, client)

		time.Sleep(time.Millisecond)
	}
}
