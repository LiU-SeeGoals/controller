package ai

import (
	// "fmt"
	"fmt"
	"math"

	"github.com/LiU-SeeGoals/controller/internal/info"
)

type MovementComposition struct{}

const ( // AvoidCollision constants
	detectionRadius = float32(500.0) // How close is "too close" to another robot (mm)
	baseDamping     = float32(50.0)     // baseline damping
	brakeDampingMax = float32(5.0)    // max damping when very close
	wGyroMax        = float32(1300.0) // max gyroscopic gain
	stepClamp       = float32(2000.0) // maximum allowed step (mm)
)

const ( // SimpleAvoid constants
	kAttraction  = float32(1.0)   // how strongly we pull toward the goal
	kRepulsion   = float32(10.0) // how strongly we push away from nearby obstacles
	detectionRad = float32(500.0) // detection radius for obstacles (mm)
)

// SimpleAvoid: Just a force pulling the robot toward a goal,
// and a repulsive force pushing away from obstacles within detectionRad.
func (mc *MovementComposition) SimpleAvoid(robot *info.Robot, goal info.Position, gs *info.GameState) info.Position {
	// Current state
	obstacles := getObstacles(gs, robot)
	pos := robot.GetPosition()

	// 1) Attractive force
	Fatt := computeAttractiveForce(pos, goal, kAttraction)

	// 2) Repulsive force from obstacles
	Frep := computeRepulsiveForce(pos, obstacles, kRepulsion, detectionRad)

	// 3) Sum the forces => acceleration (mass = 1)
	Ftotal := Fatt.Add(&Frep)

	newPos := pos.Add(&Ftotal)

	return newPos
}

// computeAttractiveForce: -k*(pos - goal), i.e. a vector from pos→goal
func computeAttractiveForce(pos, goal info.Position, k float32) info.Position {
	dx := goal.X - pos.X
	dy := goal.Y - pos.Y
	return info.Position{
		X: k * dx,
		Y: k * dy,
	}
}

func computeRepulsiveForce(pos info.Position, obstacles []info.Position, k, rDet float32) info.Position {
	force := info.Position{X: 0, Y: 0}
	for _, obs := range obstacles {
		dx := pos.X - obs.X
		dy := pos.Y - obs.Y

		dist := float32(math.Sqrt(float64(dx*dx + dy*dy)))
		fmt.Println("Distance: ", dist)

		if dist < rDet && dist != 0{

			// push away
			mag := k * (rDet - dist)
			// direction: from obs → pos
			force.X += mag * (dx / dist)
			force.Y += mag * (dy / dist)
		}
	}
	return force
}

func getObstacles(gs *info.GameState, robot *info.Robot) []info.Position {

	var obstacles []info.Position

	// Handle own team avoiding
	for _, otherRobot := range gs.Yellow_team {
		// Avoid self
		if otherRobot.GetID() != robot.GetID() {
			obstacles = append(obstacles, otherRobot.GetPosition())

		}
	}

	for _, otherRobot := range gs.Blue_team {
		obstacles = append(obstacles, otherRobot.GetPosition())
	}

	return obstacles
}

// AvoidCollision computes a single-step “waypoint” that points the robot
// in the direction of the sum of potential, gyroscopic, and damping/braking forces.
func (mc *MovementComposition) AvoidCollision(robot *info.Robot, goal info.Position, gs *info.GameState) info.Position {
	// Current position and velocity
	robotPos := robot.GetPosition()
	oldVel := robot.GetVelocity()

	// 1) Find nearest obstacle/robot
	nearestRobot := getNearestRobot(gs, robot)

	// 2) Potential force (robot → goal)
	Fp := computePotentialForce(robotPos, goal)

	// 3) Gyroscopic force
	Fg := computeGyroscopicForce(robotPos, oldVel, nearestRobot)

	// 4) Damping + Braking force
	Fd := computeDampingForce(oldVel, robotPos, nearestRobot)

	// Sum up the three terms: total “direction” we want to move
	totalForce := Fp.Add(&Fg)
	totalForce = totalForce.Add(&Fd)

	// (Optional) Clamp the step size so we don’t send a giant jump
	stepSize := totalForce.Norm()
	if stepSize > stepClamp {
		totalForce = totalForce.Scale(stepClamp / stepSize)
	}

	// This code returns a single step from current position
	// in the direction of totalForce
	newPos := robotPos.Add(&totalForce)

	// Return newPos as the waypoint
	return newPos
}

// computePotentialForce: points from robot to goal, optionally magnitude-clamped
func computePotentialForce(robotPos, goal info.Position) info.Position {
	attractive := goal.Sub(&robotPos)
	mag := attractive.Norm()
	maxMag := float32(2000.0)
	if mag > maxMag {
		attractive = attractive.Scale(maxMag / mag)
	}
	return attractive
}

// computeGyroscopicForce: rotates velocity ±90° depending on (obsVec × vel),
// scaled by how close the other robot is and how large the velocity is.
func computeGyroscopicForce(robotPos, vel info.Position, other *info.Robot) info.Position {
	if other == nil {
		return info.Position{X: 0, Y: 0}
	}

	// Vector from us to nearest robot
	otherPos := other.GetPosition()
	obsVec := otherPos.Sub(&robotPos)

	// Speed
	speed := vel.Norm()
	if speed < 1e-8 {
		return info.Position{X: 0, Y: 0}
	}

	// Cross sign
	crossVal := obsVec.Cross2D(&vel) // 2D cross => obsVec.X*vel.Y - obsVec.Y*vel.X
	signDir := float32(1.0)
	if crossVal < 0 {
		signDir = -1.0
	}

	// ±90 deg rotation of velocity
	vPerp := info.Position{
		X: -vel.Y * signDir,
		Y: vel.X * signDir,
	}

	// Distance-based scale => stronger if close
	dist := obsVec.Norm()
	if dist < 1e-3 {
		dist = 1e-3
	}
	scale := float32(0.0)
	if dist < detectionRadius {
		scale = (detectionRadius - dist) / detectionRadius
	}
	magnitude := wGyroMax * scale * speed

	// Normalize vPerp => multiply by magnitude
	perpLen := vPerp.Norm()
	if perpLen < 1e-8 {
		return info.Position{X: 0, Y: 0}
	}
	return vPerp.Scale(magnitude / perpLen)
}

// computeDampingForce: base damping plus extra “brake” if the other robot is close
func computeDampingForce(vel, robotPos info.Position, other *info.Robot) info.Position {
	if other == nil {
		return vel.Scale(-baseDamping)
	}

	otherPos := other.GetPosition()
	distance := robotPos.Distance(&otherPos)

	if distance >= detectionRadius {
		// beyond detection => base damping only
		return vel.Scale(-baseDamping)
	}

	ratio := (detectionRadius - distance) / detectionRadius
	D := baseDamping + ratio*(brakeDampingMax-baseDamping)
	return vel.Scale(-D)
}

// getNearestRobot scans both teams for the closest robot under detectionRadius, ignoring self.
func getNearestRobot(gs *info.GameState, robot *info.Robot) *info.Robot {
	minDist := float32(math.MaxFloat32)
	var nearest *info.Robot

	robotPos := robot.GetPosition()

	// check same team
	for _, other := range gs.GetTeam(robot.GetTeam()) {
		if other.GetID() == robot.GetID() { continue }
		if !other.IsActive() { continue }

		otherPos := other.GetPosition()
		dist := robotPos.Distance(&otherPos)
		if dist < minDist && dist < detectionRadius {
			minDist = dist
			nearest = other
		}
	}
	// check other team
	for _, other := range gs.GetOtherTeam(robot.GetTeam()) {
		if !other.IsActive() { continue }

		otherPos := other.GetPosition()
		dist := robotPos.Distance(&otherPos)
		if dist < minDist && dist < detectionRadius {
			minDist = dist
			nearest = other
		}
	}
	return nearest
}

// -----------------------------
// Note: Not integrated with new code structure
// -----------------------------

// import (
// 	"encoding/json"
// 	"fmt"
// 	"math"
// 	"net"

// 	"github.com/LiU-SeeGoals/controller/internal/info"
// 	"gonum.org/v1/gonum/mat"
// )

// const (
// 	fieldLength = 9000 // Length of the field in mm
// 	fieldWidth  = 6000 // Width of the field in mm
// 	cellSize    = 200  // Size of each cell in mm
// 	kAtt        = 1.0  // Attractive potential constant
// 	kRep        = 25.0 // Repulsive potential constant
// 	d0          = 5.0  // Distance at which repulsive potential is 0
// 	localSize   = 5    // Size of the local neighborhood
// )

// // robotPos is the current position of the robot
// // goal is the destination position
// // gs is the current game state, contains the positions of all robots ie. the obstacles
// func avoidObstacles(robot *info.Robot, goal info.Position, gs info.GameState) info.Position {

// 	// Matrix to hold the potential in the local neighborhood
// 	localGrid := mat.NewDense(localSize, localSize, nil)

// 	obstacles := getObstacles(gs, robot.GetID())
// 	obstacles = addWallObstacles(obstacles)

// 	localGrid.Apply(func(i, j int, v float64) float64 {
// 		centerOffset := int(math.Floor(localSize / 2))
// 		x := robot.GetPosition().X/cellSize + float32(i-centerOffset)
// 		y := robot.GetPosition().Y/cellSize + float32(j-centerOffset)

// 		// Compute the attractive potential
// 		attractive := computeAttractivePotential(x, y, goal.X/cellSize, goal.Y/cellSize)

// 		// Compute the repulsive potential
// 		repulsive := computeRepulsivePotential(x, y, obstacles, d0, kRep)
// 		// repulsive = 0.0

// 		return attractive + repulsive
// 	}, localGrid)

// 	// Send the local grid to the Python script
// 	// sendLocalGrid(localGrid)

// 	// minPotentialRow, minPotentialCol, _ := argmin(localGrid)
// 	minPotentialRow, minPotentialCol, _ := argminNeighbors(localGrid, int(math.Floor(localSize/2)), int(math.Floor(localSize/2)))

// 	// Calculate the offsets relative to the robot’s current position
// 	centerOffset := int(math.Floor(localSize / 2))

// 	offsetX := float32(minPotentialRow-centerOffset) * cellSize
// 	offsetY := float32(minPotentialCol-centerOffset) * cellSize

// 	// Apply the offsets to the robot’s current position to get the new destination
// 	newX := robot.GetPosition().X + offsetX
// 	newY := robot.GetPosition().Y + offsetY

// 	return info.Position{X: newX, Y: newY}
// }

// func addWallObstacles(obstacles []info.Position) []info.Position {

// 	// Add the walls as obstacles
// 	padding := 800
// 	halfFieldWidth := (fieldWidth + padding) / 2
// 	halfFieldLength := (fieldLength + padding) / 2
// 	robotRadius := 50
// 	for x := -halfFieldLength - robotRadius; x < halfFieldLength+robotRadius; x += 2 * robotRadius {
// 		obstacles = append(obstacles, info.Position{X: float32(x), Y: float32(halfFieldWidth)})
// 		obstacles = append(obstacles, info.Position{X: float32(x), Y: float32(-halfFieldWidth)})
// 	}
// 	for y := -halfFieldWidth - robotRadius; y < halfFieldWidth+robotRadius; y += 2 * robotRadius {
// 		obstacles = append(obstacles, info.Position{X: float32(halfFieldLength), Y: float32(y)})
// 		obstacles = append(obstacles, info.Position{X: float32(-halfFieldLength), Y: float32(y)})
// 	}
// 	return obstacles

// }


// func computeAttractivePotential(x, y, goalX, goalY float32) float64 {
// 	dx := float64(x - goalX)
// 	dy := float64(y - goalY)
// 	return 0.5 * kAtt * math.Sqrt(math.Pow(dx, 2)+math.Pow(dy, 2))
// }

// // calculateRepulsivePotential calculates the repulsive potential from obstacles
// func computeRepulsivePotential(x, y float32, obstacles []info.Position, d0, kRep float64) float64 {
// 	repulsive := 0.0
// 	for _, obstacle := range obstacles {
// 		obstacleX, obstacleY := obstacle.X/cellSize, obstacle.Y/cellSize

// 		dx := float64(x - obstacleX)
// 		dy := float64(y - obstacleY)

// 		distance := math.Sqrt(dx*dx + dy*dy)
// 		if distance < 2 {
// 			repulsive += 100
// 		} else if distance < d0 && distance != 0 {
// 			repulsive += 0.5 * kRep * math.Pow((1/distance)-(1/d0), 2)
// 		}
// 	}
// 	return repulsive
// }

// func argminNeighbors(m *mat.Dense, row, col int) (int, int, float64) {
// 	minValue := math.MaxFloat64
// 	minRow, minCol := -1, -1

// 	rows, cols := m.Dims()

// 	// Define the relative positions of the 8 neighbors
// 	directions := [][2]int{
// 		{-1, -1}, {-1, 0}, {-1, 1}, // Top-left, Top, Top-right
// 		{0, -1}, {0, 1}, // Left,       Right
// 		{1, -1}, {1, 0}, {1, 1}, // Bottom-left, Bottom, Bottom-right
// 	}

// 	for _, d := range directions {
// 		neighborRow := row + d[0]
// 		neighborCol := col + d[1]

// 		// Check bounds
// 		if neighborRow >= 0 && neighborRow < rows && neighborCol >= 0 && neighborCol < cols {
// 			val := m.At(neighborRow, neighborCol)
// 			if val < minValue {
// 				minValue = val
// 				minRow, minCol = neighborRow, neighborCol
// 			}
// 		}
// 	}

// 	return minRow, minCol, minValue
// }

// func argmin(m *mat.Dense) (int, int, float64) {
// 	minValue := math.MaxFloat64
// 	minRow, minCol := 0, 0

// 	rows, cols := m.Dims()
// 	for i := 0; i < rows; i++ {
// 		for j := 0; j < cols; j++ {
// 			val := m.At(i, j)
// 			if val < minValue {
// 				minValue = val
// 				minRow, minCol = i, j
// 			}
// 		}
// 	}
// 	return minRow, minCol, minValue
// }

// var listener net.Listener
// var conn net.Conn

// func init() {
// 	// Start the server to listen for incoming connections
// 	var err error
// 	listener, err = net.Listen("tcp", ":5000") // Bind to port 5000
// 	if err != nil {
// 		fmt.Printf("Error starting server: %v\n", err)
// 		return
// 	}

// 	// Accept connections in a goroutine to allow non-blocking operations
// 	go func() {
// 		for {
// 			fmt.Println("Waiting for Python client to connect...")
// 			conn, err = listener.Accept()
// 			if err != nil {
// 				fmt.Printf("Error accepting connection: %v\n", err)
// 				continue
// 			}
// 			fmt.Println("Python client connected.")
// 		}
// 	}()
// }

// func sendLocalGrid(localGrid *mat.Dense) {
// 	if conn == nil {
// 		fmt.Println("No active connection to send data")
// 		return
// 	}

// 	rows, cols := localGrid.Dims()
// 	data := make([][]float64, rows)
// 	for i := 0; i < rows; i++ {
// 		row := make([]float64, cols)
// 		for j := 0; j < cols; j++ {
// 			row[j] = localGrid.At(i, j)
// 		}
// 		data[i] = row
// 	}

// 	// Serialize to JSON
// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		fmt.Printf("Error serializing localGrid: %v\n", err)
// 		return
// 	}

// 	// Send the JSON data followed by a newline
// 	_, err = conn.Write(append(jsonData, '\n')) // Append newline for framing
// 	if err != nil {
// 		fmt.Printf("Error sending localGrid: %v\n", err)
// 		conn = nil // Reset connection if sending fails
// 	}
// }
