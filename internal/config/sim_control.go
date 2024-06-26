package config

import (
	"fmt"

	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/helper"
	"github.com/LiU-SeeGoals/proto_go/simulation"
)

// The simulator have a lot of things that can be configured.
// This configuration is done with proto messages on port 10300 (not the port for teams).
type simControl struct {
	client *client.SimClient
}

func NewSimControl() *simControl {
	simClient := client.NewSimClient(GetSimControlAddress())
	simClient.Init()
	return &simControl{
		client: simClient,
	}
}

func (sc *simControl) TurnOffCameraRealism() {
	fmt.Println("Not yet implemented")
}

func (sc *simControl) TurnOnCameraRealism() {
	fmt.Println("Not yet implemented")
}

func (sc *simControl) SetPresentRobots(presentYellow []int, presentBlue []int) {
	TOTAL_ROBOTS := 11
	x := float32(1.0)           // X-coordinate
	y := float32(1.0)           // Y-coordinate
	orientation := float32(0.0) // Approx. 45 degrees in radians
	vx := float32(0.0)          // Velocity towards x-axis
	vy := float32(0.0)          // Velocity towards y-axis
	vAngular := float32(0.0)    // Angular velocity

	robotList := []*simulation.TeleportRobot{}

	for i := 0; i < TOTAL_ROBOTS; i++ {
		present := false
		if helper.Contains(presentYellow, i) { // all robots after the number we want --> set to not present
			present = true
		}

		idNum := uint32(i)
		team := simulation.Team_YELLOW
		id := simulation.SimRobotId{
			Id:   &idNum,
			Team: &team,
		}

		teleportRobot := &simulation.TeleportRobot{
			Id:          &id,
			X:           &x,
			Y:           &y,
			Orientation: &orientation,
			VX:          &vx,
			VY:          &vy,
			VAngular:    &vAngular,
			Present:     &present,
		}

		robotList = append(robotList, teleportRobot)
	}

	for i := 0; i < TOTAL_ROBOTS; i++ {
		present := false
		if helper.Contains(presentBlue, i) { // all robots after the number we want --> set to not present
			present = true
		}

		idNum := uint32(i)
		team := simulation.Team_BLUE
		id := simulation.SimRobotId{
			Id:   &idNum,
			Team: &team,
		}

		teleportRobot := &simulation.TeleportRobot{
			Id:          &id,
			X:           &x,
			Y:           &y,
			Orientation: &orientation,
			VX:          &vx,
			VY:          &vy,
			VAngular:    &vAngular,
			Present:     &present,
		}

		robotList = append(robotList, teleportRobot)
	}

	simControl := &simulation.SimulatorControl{
		TeleportRobot:   robotList,
		TeleportBall:    nil,
		SimulationSpeed: nil,
	}

	simCommand := &simulation.SimulatorCommand{
		Control: simControl,
		Config:  nil,
	}

	sc.client.Send(simCommand)
}

func (sc *simControl) SetRobotDimentions() {
	fmt.Println("Not yet implemented")
}

func (sc *simControl) RobotStartPositionConfig1(numberOfRobots int) {
	fmt.Println("Not yet implemented")
}

func (sc *simControl) RobotStartPositionConfig2(numberOfRobots int) {
	fmt.Println("Not yet implemented")
}
