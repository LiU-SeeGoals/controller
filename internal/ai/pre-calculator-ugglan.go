package ai

import (
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"gonum.org/v1/gonum/mat"
)

const NUM_ROWS = 3
const NUM_COLS = 6
const NUM_CHANNELS = 3

type PreCalculator struct {
	analysis GameAnalysis
}

// Constructor for the PreCalculator
func NewPreCalculator(field gamestate.Field) *PreCalculator {
	pc := &PreCalculator{
		analysis: *newAnalysis(field.FieldLengt, field.FieldWidth),
	}
	return pc
}

// GameAnalysis constructor
func newAnalysis(fieldLength, fieldWidth int32) *GameAnalysis {
	analysis := GameAnalysis{}
	zones := [NUM_ROWS * NUM_COLS]Zone{}

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
		robots:             nil,
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

// The pitch is divied into 18 zones, numbered 0-17, starting from the defensive
// end to the attacking end, and from (goalkeepers perspective) left to right across the field.
type Zone struct {
	id                 int                // id of the zone, 0-17
	robots             []*gamestate.Robot // robots in the zone
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
	zones         [NUM_ROWS * NUM_COLS]Zone // The pitch is divided into 18 zones, 3 rows and 6 columns
	channels      [NUM_CHANNELS]Channel     // The pitch is divided into 3 channels
	inPossession  bool                      // true if the team is in possession of the ball
	zoneLength    float32                   // length of each zoneLength
	zoneWidth     float32                   // width of each zoneWidth
	channelWidth  float32                   // width of each channelWidth
	channelLength float32                   // length of each channelLength
}

func (an *GameAnalysis) updateZones(gamestateObj *gamestate.GameState) {

	// Reset the zones
	for i := 0; i < NUM_ROWS*NUM_COLS; i++ {
		an.zones[i].robots = nil
	}
	blueCount := [NUM_ROWS * NUM_COLS]int{}
	yellowCount := [NUM_ROWS * NUM_COLS]int{}

	// count blue robots in each zone
	for _, robot := range gamestateObj.Blue_team {
		col := robot.GetPosition().AtVec(0) / float64(an.zoneLength)
		row := robot.GetPosition().AtVec(1) / float64(an.zoneWidth)
		blueCount[int(row)*NUM_COLS+int(col)]++
	}

	// count yellow robots in each zone
	for _, robot := range gamestateObj.Yellow_team {
		col := robot.GetPosition().AtVec(0) / float64(an.zoneLength)
		row := robot.GetPosition().AtVec(1) / float64(an.zoneWidth)
		yellowCount[int(row)*NUM_COLS+int(col)]++
	}

	// calculate the proportion of robots in each zone
	for i := range an.zones {
		totalRobots := blueCount[i] + yellowCount[i]
		if totalRobots > 0 {
			if blueCount[i] == 0 && yellowCount[i] != 0 {
				an.zones[i].controlProbability = 1.0
			} else {
				an.zones[i].controlProbability = float64(yellowCount[i]) / float64(totalRobots)
			}
		} else {
			an.zones[i].controlProbability = 0.5
		}
	}
}

func (pc *PreCalculator) Process(gamestateObj *gamestate.GameState) *GameAnalysis {
	pc.analysis.updateZones(gamestateObj)
	// pc.updateChannels(gamestateObj)
	// pc.updatePossession(gamestateObj)

	return &pc.analysis
}
