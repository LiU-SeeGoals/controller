package main

import (
	"fmt"

	"github.com/LiU-SeeGoals/controller/internal/proto/ssl_vision"
	"github.com/LiU-SeeGoals/controller/internal/receiver"
)

func main() {
	var packet ssl_vision.SSL_WrapperPacket
	packetChan := setupSSLVisionReceiver()

	var detect *ssl_vision.SSL_DetectionFrame = nil
	for {
		packet = <-packetChan

		detect = packet.GetDetection()

		fmt.Println("-- Blue team --")
		for _, robot := range detect.GetRobotsBlue() {
			fmt.Printf("id: %v \n", robot.GetRobotId())
			fmt.Printf("x: %v \n", robot.GetX())
			fmt.Printf("y: %v \n", robot.GetY())
		}

		fmt.Println("-- Yellow team --")
		for _, robot := range detect.GetRobotsYellow() {
			fmt.Printf("id: %v \n", robot.GetRobotId())
			fmt.Printf("x: %v \n", robot.GetX())
			fmt.Printf("y: %v \n", robot.GetY())
		}

		fmt.Println("-- Ball --")
		for _, ball := range detect.GetBalls() {
			fmt.Printf("x: %v \n", ball.GetX())
			fmt.Printf("y: %v \n", ball.GetY())
			fmt.Printf("z: %v \n", ball.GetZ())
		}
	}

}

// Start a SSL Vision receiver, returns a channel from
// which SSL wrapper packets can be obtained.
func setupSSLVisionReceiver() chan ssl_vision.SSL_WrapperPacket {
	recv := receiver.NewSSLReceiver("224.5.23.2:10020")
	recv.Connect()

	c := make(chan ssl_vision.SSL_WrapperPacket)
	go recv.Receive(c)

	return c
}
