package main

import (
	"fmt"
	"net"

	"github.com/LiU-SeeGoals/controller/internal/parsed_vision"

	"google.golang.org/protobuf/proto"
)

func main() {
	// Just change this two variable
	// There is no error handling
	// Good luck :)
	address := "127.0.0.1:1234"

	parsedFrame := parsed_vision.ParsedFrame{} 
	// Initialize Robots
	robots := []parsed_vision.Robot{
		{RobotId: 1, Pos: &parsed_vision.Vector2{X: 1.0, Y: 2.0}, Orientation: 0.5, Vel: &parsed_vision.Vector2{X: 0.1, Y: 0.2}},
		{RobotId: 2, Pos: &parsed_vision.Vector2{X: 1.0, Y: 2.0}, Orientation: 0.5, Vel: &parsed_vision.Vector2{X: 0.1, Y: 0.2}},
		{RobotId: 3, Pos: &parsed_vision.Vector2{X: 1.0, Y: 2.0}, Orientation: 0.5, Vel: &parsed_vision.Vector2{X: 0.1, Y: 0.2}},
		{RobotId: 4, Pos: &parsed_vision.Vector2{X: 1.0, Y: 2.0}, Orientation: 0.5, Vel: &parsed_vision.Vector2{X: 0.1, Y: 0.2}},
		{RobotId: 5, Pos: &parsed_vision.Vector2{X: 1.0, Y: 2.0}, Orientation: 0.5, Vel: &parsed_vision.Vector2{X: 0.1, Y: 0.2}},
		{RobotId: 6, Pos: &parsed_vision.Vector2{X: 1.0, Y: 2.0}, Orientation: 0.5, Vel: &parsed_vision.Vector2{X: 0.1, Y: 0.2}},
	}

	// Adding robots to parsedFrame
	for _, robot := range robots {
		parsedFrame.Robots = append(parsedFrame.Robots, &robot)
	}

	// Initialize Ball
	ball := parsed_vision.Vector2{X: 3.0, Y: 4.0}

	// Adding ball to parsedFrame
	parsedFrame.Ball = &ball

	var err error = nil
	connection, _ := net.Dial("udp", address)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
	}

	bytes, _ := proto.Marshal(&parsedFrame)
	connection.Write(bytes)
}
