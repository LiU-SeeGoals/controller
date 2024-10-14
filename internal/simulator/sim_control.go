package simulator

import (
	"fmt"

	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/config"
	"github.com/LiU-SeeGoals/controller/internal/helper"
	"github.com/LiU-SeeGoals/proto_go/simulation"
)

// The simulator have a lot of things that can be configured.
// This configuration is done with proto messages on port 10300 (not the port for teams).
type simControl struct {
	client *client.SimClient
}

func NewSimControl() *simControl {
	simClient := client.NewSimClient(config.GetSimControlAddress())
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
  generateCoordinates := func(x, min_y, max_y float32) [][2]float32 {
    coords := make([][2]float32, numberOfRobots)
    step := (max_y - min_y) / float32(numberOfRobots - 1)
    for i := 0; i < numberOfRobots; i++ {
      y := min_y + step * float32(i)
      coords[i] = [2]float32{x, y}
    }
    return coords
  }

  blueCoords := generateCoordinates(1, -2, 2)
  yellowCoords := generateCoordinates(-1, -2, 2)

  for robot_id := 0; robot_id < numberOfRobots; robot_id++ {
    x_blue := blueCoords[robot_id][0]
    y_blue := blueCoords[robot_id][1]
    id := uint32(robot_id)
    team := simulation.Team_BLUE
    sc.TeleportRobot(x_blue, y_blue, id, team)

    x_yellow := yellowCoords[robot_id][0]
    y_yellow := yellowCoords[robot_id][1]
    id = uint32(robot_id)
    team = simulation.Team_YELLOW
    sc.TeleportRobot(x_yellow, y_yellow, id, team)
  }

}

func (sc *simControl) RobotStartPositionConfig2(numberOfRobots int) {
	fmt.Println("Not yet implemented")
}

func (sc *simControl) TeleportRobot(x float32, y float32, id uint32, team simulation.Team) {
	fmt.Println(x, y)
	// Set default values for orientation and velocities
	orientation := float32(0.0) // Approx. 45 degrees in radians
	vx := float32(0.0)          // Velocity towards x-axis
	vy := float32(0.0)          // Velocity towards y-axis
	vAngular := float32(0.0)    // Angular velocity
	present := true             // Teleport indicates the robot is present

	// Create the robot ID structure
	robotId := simulation.SimRobotId{
		Id:   &id,
		Team: &team,
	}

	// Create the TeleportRobot structure with the new position and parameters
	teleportRobot := &simulation.TeleportRobot{
		Id:          &robotId,
		X:           &x,
		Y:           &y,
		Orientation: &orientation,
		VX:          &vx,
		VY:          &vy,
		VAngular:    &vAngular,
		Present:     &present,
	}

	// Prepare the command with the single robot teleportation
	simControl := &simulation.SimulatorControl{
		TeleportRobot:   []*simulation.TeleportRobot{teleportRobot},
		TeleportBall:    nil,
		SimulationSpeed: nil,
	}

	simCommand := &simulation.SimulatorCommand{
		Control: simControl,
		Config:  nil,
	}

	// Send the command to teleport the robot
	sc.client.Send(simCommand)
}
