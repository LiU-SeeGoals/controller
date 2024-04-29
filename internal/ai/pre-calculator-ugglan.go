package ai

import (
	"math"

	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"gonum.org/v1/gonum/mat"
)

const NUM_ROWS = 3
const NUM_COLS = 3
const NUM_CHANNELS = 3

type PreCalculator struct {
	analysis *GameAnalysis
}

// Constructor for the PreCalculator
func NewPreCalculator(fieldLength, fieldWidth int32) *PreCalculator {
	pc := &PreCalculator{
		analysis: newAnalysis(fieldLength, fieldWidth),
	}
	return pc
}

// GameAnalysis constructor
func newAnalysis(fieldLength, fieldWidth int32) *GameAnalysis {
	analysis := GameAnalysis{}
	zones := [NUM_ROWS * NUM_COLS]Zone{}
	analysis.zoneLength = float32(fieldLength) / float32(NUM_COLS)
	analysis.zoneWidth = float32(fieldWidth) / float32(NUM_ROWS)

	// Initialize the zones
	for i := 0; i < NUM_ROWS; i++ {
		for j := 0; j < NUM_COLS; j++ {
			x := float64(j) * float64(fieldLength) / float64(NUM_COLS)
			y := float64(i) * float64(fieldWidth) / float64(NUM_ROWS)
			id := i*NUM_COLS + j
			zones[id] = *newZone(id, *mat.NewVecDense(2, []float64{x, y}))
		}
	}
	analysis.zones = zones

	return &analysis
}

// Helper function to get the adjacent zones of a given zone
func adjacentZones(id int) [4]int {
	directions := [4]int{-1, 1, -NUM_COLS, NUM_COLS}
	adjacentZones := [4]int{}
	for i, direction := range directions {
		// Check if the adjacent zone is in of bounds
		if id+direction >= 0 || id+direction < NUM_ROWS*NUM_COLS {
			adjacentZones[i] = id + direction
		}
	}
	return adjacentZones
}

// Zone Constructor
func newZone(id int, centerCoordinates mat.VecDense) *Zone {
	zone := Zone{
		id:                 id,
		yellow_robots:      nil,
		blue_robots:        nil,
		controlProbability: 0.0,
		centerCoordinates:  centerCoordinates,
		adjacentZones:      adjacentZones(id),
	}
	return &zone
}

// Channel constructor
func newChannel(id int, associatedZones [NUM_COLS]int) *Channel {
	channel := Channel{
		id:              id,
		AssociatedZones: associatedZones,
		robots:          nil,
	}
	return &channel
}

// The pitch is divied into 9 zones, numbered 0-8, starting from the defensive
// end to the attacking end, and from (goalkeepers perspective) left to right across the field.
type Zone struct {
	id                 int                // id of the zone, 0-8
	blue_robots        []*gamestate.Robot // blue robots in the zone
	yellow_robots      []*gamestate.Robot // yellow robots in the zone
	controlProbability float64            // probability of maintaining control of the ball in this zone
	centerCoordinates  mat.VecDense       // coordinates of the center of the zone
	adjacentZones      [4]int             // array of adjacent zones ids
	width              float32            // width of the Zone
	length             float32            // height of the Zone
}

// The pitch is divided into 3 channels, left wing(0), center channel(1), and right wing(2).
type Channel struct {
	id                 int                // id of channel. 0: left wing, 1: center channel, 2: right wing
	AssociatedZones    [NUM_COLS]int      // IDs of zones that fall within this channel
	robots             []*gamestate.Robot // robots in the channel
	controlProbability float64            // probability of maintaining control of the ball in this channel
}

// Struct to hold the analysis of the gamestate
type GameAnalysis struct {
	zones         [NUM_ROWS * NUM_COLS]Zone // The pitch is divided into 9 zones, 3 rows and 3 columns
	channels      [NUM_CHANNELS]Channel     // The pitch is divided into 3 channels
	inPossession  bool                      // true if the team is in possession of the ball
	zoneLength    float32                   // length of each zoneLength
	zoneWidth     float32                   // width of each zoneWidth
	channelWidth  float32                   // width of each channelWidth
	channelLength float32                   // length of each channelLength
}

// BROKEN
func (an *GameAnalysis) updateZones(gamestateObj *gamestate.GameState) {

	// Reset the zones
	for i := 0; i < NUM_ROWS*NUM_COLS; i++ {
		an.zones[i].yellow_robots = []*gamestate.Robot{}
		an.zones[i].blue_robots = []*gamestate.Robot{}
		an.zones[i].controlProbability = 0.49
	}
	// count blue robots in each zone and add them to the zone
	for _, robot := range gamestateObj.Blue_team {
		col := math.Floor(robot.GetPosition().AtVec(0) / float64(an.zoneLength))
		row := math.Floor(robot.GetPosition().AtVec(1) / float64(an.zoneWidth))
		an.zones[int(row)*NUM_COLS+int(col)].blue_robots = append(an.zones[int(row)*NUM_COLS+int(col)].blue_robots, robot)
	}

	// count yellow robots in each zone
	for _, robot := range gamestateObj.Yellow_team {
		col := math.Floor(robot.GetPosition().AtVec(0) / float64(an.zoneLength))
		row := math.Floor(robot.GetPosition().AtVec(1) / float64(an.zoneWidth))
		an.zones[int(row)*NUM_COLS+int(col)].yellow_robots = append(an.zones[int(row)*NUM_COLS+int(col)].yellow_robots, robot)
	}

	// calculate the proportion of robots in each zone
	for i := range an.zones {
		blueCount := len(an.zones[i].blue_robots)
		yellowCount := len(an.zones[i].yellow_robots)
		totalRobots := blueCount + yellowCount
		if totalRobots > 0 {
			if blueCount == 0 && yellowCount != 0 {
				an.zones[i].controlProbability = 1.0
			} else {
				an.zones[i].controlProbability = float64(yellowCount) / float64(totalRobots)
			}
		} else {
			an.zones[i].controlProbability = 0.49
		}
	}
}

func (pc *PreCalculator) Process(gamestateObj *gamestate.GameState) *GameAnalysis {
	// pc.analysis.updateZones(gamestateObj)
	// pc.updateChannels(gamestateObj)
	// pc.updatePossession(gamestateObj)

	return pc.analysis
}
