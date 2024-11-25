package ai

import (
	"encoding/json"
	"fmt"
	"math"
	"net"

	"github.com/LiU-SeeGoals/controller/internal/state"
	"gonum.org/v1/gonum/mat"
)

const (
	fieldLength = 9000 // Length of the field in mm
	fieldWidth  = 6000 // Width of the field in mm
	cellSize    = 200  // Size of each cell in mm
	kAtt        = 1.0  // Attractive potential constant
	kRep        = 25.0 // Repulsive potential constant
	d0          = 5.0  // Distance at which repulsive potential is 0
	localSize   = 5    // Size of the local neighborhood
)

// robotPos is the current position of the robot
// goal is the destination position
// gs is the current game state, contains the positions of all robots ie. the obstacles
func getObstacleFreeDest(robot *state.Robot, goal state.Position, gs state.GameState) state.Position {

	// Matrix to hold the potential in the local neighborhood
	localGrid := mat.NewDense(localSize, localSize, nil)

	obstacles := getObstacles(gs, robot.GetID())
	obstacles = addWallObstacles(obstacles)

	localGrid.Apply(func(i, j int, v float64) float64 {
		centerOffset := int(math.Floor(localSize / 2))
		x := robot.GetPosition().X/cellSize + float32(i-centerOffset)
		y := robot.GetPosition().Y/cellSize + float32(j-centerOffset)

		// Compute the attractive potential
		attractive := computeAttractivePotential(x, y, goal.X/cellSize, goal.Y/cellSize)

		// Compute the repulsive potential
		repulsive := computeRepulsivePotential(x, y, obstacles, d0, kRep)
		// repulsive = 0.0

		return attractive + repulsive
	}, localGrid)

	// Send the local grid to the Python script
	// sendLocalGrid(localGrid)

	// minPotentialRow, minPotentialCol, _ := argmin(localGrid)
	minPotentialRow, minPotentialCol, _ := argminNeighbors(localGrid, int(math.Floor(localSize/2)), int(math.Floor(localSize/2)))

	// Calculate the offsets relative to the robot’s current position
	centerOffset := int(math.Floor(localSize / 2))

	offsetX := float32(minPotentialRow-centerOffset) * cellSize
	offsetY := float32(minPotentialCol-centerOffset) * cellSize

	// Apply the offsets to the robot’s current position to get the new destination
	newX := robot.GetPosition().X + offsetX
	newY := robot.GetPosition().Y + offsetY

	return state.Position{X: newX, Y: newY}
}

func addWallObstacles(obstacles []state.Position) []state.Position {

	// Add the walls as obstacles
	padding := 800
	halfFieldWidth := (fieldWidth + padding) / 2
	halfFieldLength := (fieldLength + padding) / 2
	robotRadius := 50
	for x := -halfFieldLength - robotRadius; x < halfFieldLength+robotRadius; x += 2 * robotRadius {
		obstacles = append(obstacles, state.Position{X: float32(x), Y: float32(halfFieldWidth)})
		obstacles = append(obstacles, state.Position{X: float32(x), Y: float32(-halfFieldWidth)})
	}
	for y := -halfFieldWidth - robotRadius; y < halfFieldWidth+robotRadius; y += 2 * robotRadius {
		obstacles = append(obstacles, state.Position{X: float32(halfFieldLength), Y: float32(y)})
		obstacles = append(obstacles, state.Position{X: float32(-halfFieldLength), Y: float32(y)})
	}
	return obstacles

}

func getObstacles(gs state.GameState, id state.ID) []state.Position {

	var obstacles []state.Position
	for _, robot := range gs.Yellow_team {
		obstacles = append(obstacles, robot.GetPosition())
	}

	// Handle own team avoiding
	for _, robot := range gs.Blue_team {
		// Avoid self
		if robot.GetID() != id {
			obstacles = append(obstacles, robot.GetPosition())
		}
	}
	return obstacles
}

func computeAttractivePotential(x, y, goalX, goalY float32) float64 {
	dx := float64(x - goalX)
	dy := float64(y - goalY)
	return 0.5 * kAtt * math.Sqrt(math.Pow(dx, 2)+math.Pow(dy, 2))
}

// calculateRepulsivePotential calculates the repulsive potential from obstacles
func computeRepulsivePotential(x, y float32, obstacles []state.Position, d0, kRep float64) float64 {
	repulsive := 0.0
	fmt.Println("nr of obstacles", len(obstacles))
	for _, obstacle := range obstacles {
		fmt.Println("Obstacle\n", obstacle)
		obstacleX, obstacleY := obstacle.X/cellSize, obstacle.Y/cellSize
		fmt.Println("Obstacle", obstacleX, obstacleY)
		fmt.Println("Robot", x, y)

		dx := float64(x - obstacleX)
		dy := float64(y - obstacleY)

		distance := math.Sqrt(dx*dx + dy*dy)
		if distance < 2 {
			repulsive += 100
		} else if distance < d0 && distance != 0 {
			repulsive += 0.5 * kRep * math.Pow((1/distance)-(1/d0), 2)
		}
	}
	return repulsive
}

func argminNeighbors(m *mat.Dense, row, col int) (int, int, float64) {
	minValue := math.MaxFloat64
	minRow, minCol := -1, -1

	rows, cols := m.Dims()

	// Define the relative positions of the 8 neighbors
	directions := [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1}, // Top-left, Top, Top-right
		{0, -1}, {0, 1}, // Left,       Right
		{1, -1}, {1, 0}, {1, 1}, // Bottom-left, Bottom, Bottom-right
	}

	for _, d := range directions {
		neighborRow := row + d[0]
		neighborCol := col + d[1]

		// Check bounds
		if neighborRow >= 0 && neighborRow < rows && neighborCol >= 0 && neighborCol < cols {
			val := m.At(neighborRow, neighborCol)
			if val < minValue {
				minValue = val
				minRow, minCol = neighborRow, neighborCol
			}
		}
	}

	return minRow, minCol, minValue
}

func argmin(m *mat.Dense) (int, int, float64) {
	minValue := math.MaxFloat64
	minRow, minCol := 0, 0

	rows, cols := m.Dims()
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			val := m.At(i, j)
			if val < minValue {
				minValue = val
				minRow, minCol = i, j
			}
		}
	}
	return minRow, minCol, minValue
}

var listener net.Listener
var conn net.Conn

func init() {
	// Start the server to listen for incoming connections
	var err error
	listener, err = net.Listen("tcp", ":5000") // Bind to port 5000
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}

	// Accept connections in a goroutine to allow non-blocking operations
	go func() {
		for {
			fmt.Println("Waiting for Python client to connect...")
			conn, err = listener.Accept()
			if err != nil {
				fmt.Printf("Error accepting connection: %v\n", err)
				continue
			}
			fmt.Println("Python client connected.")
		}
	}()
}

func sendLocalGrid(localGrid *mat.Dense) {
	if conn == nil {
		fmt.Println("No active connection to send data")
		return
	}

	rows, cols := localGrid.Dims()
	data := make([][]float64, rows)
	for i := 0; i < rows; i++ {
		row := make([]float64, cols)
		for j := 0; j < cols; j++ {
			row[j] = localGrid.At(i, j)
		}
		data[i] = row
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error serializing localGrid: %v\n", err)
		return
	}

	// Send the JSON data followed by a newline
	_, err = conn.Write(append(jsonData, '\n')) // Append newline for framing
	if err != nil {
		fmt.Printf("Error sending localGrid: %v\n", err)
		conn = nil // Reset connection if sending fails
	}
}
