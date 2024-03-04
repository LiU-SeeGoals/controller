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
func NewPreCalculator() *PreCalculator {
	// zones := make([][]float64, NUM_ROWS)
	// for i := range zones {
	// 	zones[i] = make([]float64, NUM_COLS)
	// }

	// channels := make([]float64, NUM_CHANNELS)

	// analysis := GameAnalysis{}
	pc := &PreCalculator{
		analysis: *newAnalysis(),
	}
	return pc
}

// The pitch is divied into 18 zones, 0-17, starting from the defensive
// end to the attacking end, and from (goalkeepers perspective) left to right across the field.
type Zone struct {
	id                 int                // id of the zone, 0-17
	robots             []*gamestate.Robot // robots in the zone
	controlProbability float64            // probability of maintaining or gaining control of the ball in this zone
	centerCoordinates  mat.VecDense         // coordinates of the center of the zone
	adjacentZones      []int             // slice of adjecent zones
}

type Channel struct {
	id              int                // id of channel. 0: left wing, 1: center channel, 2: right wing
	AssociatedZones [NUM_COLS]int      // IDs of zones that fall within this channel
	robots          []*gamestate.Robot // robots in the channel
}

// Struct to hold the analysis of the gamestate
type GameAnalysis struct {
	zones        [NUM_ROWS][NUM_COLS]Zone // The pitch is divided into 18 zones, 3 rows and 6 columns
	channels     [NUM_CHANNELS]Channel    // The pitch is divided into 3 channels
	inPossession bool                     // true if the team is in possession of the ball
}

// GameAnalysis constructor
func newAnalysis() *GameAnalysis {
	analysis := GameAnalysis{}
	return &analysis
}

func (pc *PreCalculator) Process(gamestateObj *gamestate.GameState) *GameAnalysis {
	// pc.updateZones(gamestateObj)
	// pc.updateChannels(gamestateObj)
	// pc.updatePossession(gamestateObj)

	return &pc.analysis
}
