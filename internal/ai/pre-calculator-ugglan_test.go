package ai

import (
	"fmt"
	"testing"

	"github.com/LiU-SeeGoals/controller/internal/gamestate"
)

func newMockField(l, w int32) *gamestate.Field {
	field := gamestate.Field{}
	field.FieldLength = l
	field.FieldWidth = w
	return &field
}

func newMockRobot(id int, team gamestate.Team, x, y float32) *gamestate.Robot {
	robot := gamestate.NewRobot(id, team)
	return robot
}

func newMockBall() *gamestate.Ball {
	ball := gamestate.NewBall()
	return ball
}

func newMockGameState() *gamestate.GameState {
	gs := &gamestate.GameState{}

	gs.Ball = newMockBall()
	gs.Field = newMockField(6, 3)

	for i := 0; i < gamestate.TEAM_SIZE; i++ {
		gs.Blue_team[i] = newMockRobot(i, gamestate.Blue, 0, 0)
		gs.Yellow_team[i] = newMockRobot(i, gamestate.Yellow, 0, 0)
	}
	return gs
}

func addRobots(gs *gamestate.GameState) {
	for i := 0; i < gamestate.TEAM_SIZE; i++ {
		gs.Blue_team[i] = newMockRobot(i, gamestate.Blue, 1, 1)
		gs.Yellow_team[i] = newMockRobot(i, gamestate.Yellow, 1, 1)
	}
}

func TestUpdateZones(t *testing.T) {
	// Create necessary objects
	field := newMockField(2*NUM_COLS, 2*NUM_ROWS)
	gs := newMockGameState()
	addRobots(gs)
	pc := NewPreCalculator(field)

	// all robots are in the same zone 
	// so the control probability should be 0.5
	pc.analysis.updateZones(gs)
	if pc.analysis.zones[0].controlProbability != 0.5 {
		t.Errorf("Expected controlProbability of zone 0 to be 0.5, got %f", pc.analysis.zones[0].controlProbability)
	}
	
	// move one robot to zone 3
	gs.Yellow_team[0].SetPosition(0, 3, 0)
	pc.analysis.updateZones(gs)
	if pc.analysis.zones[3].controlProbability != 1.0 {
		t.Errorf("Expected controlProbability of zone 6 to be 1.0, got %f", pc.analysis.zones[6].controlProbability)
	}
}

func TestNewAnalysis(t *testing.T) {
	field := newMockField(NUM_COLS, NUM_ROWS)
	analysis := newAnalysis(field.FieldLength, field.FieldWidth)

	// TODO number of zones is should not be hardcoded
	if len(analysis.zones) != NUM_ROWS*NUM_COLS {
		t.Errorf("Expected 9 zones, got %d", len(analysis.zones))
	}
	if len(analysis.channels) != NUM_CHANNELS {
		t.Errorf("Expected 3 channels, got %d", len(analysis.channels))
	}
	if analysis.zoneLength != 1.0 {
		t.Errorf("Expected zoneLength to be 1.0, got %f", analysis.zoneLength)
	}
	if analysis.zoneWidth != 1.0 {
		t.Errorf("Expected zoneWidth to be 1.0, got %f", analysis.zoneWidth)
	}
	// if analysis.channelWidth != 1.0 {
	// 	t.Errorf("Expected channelwidth to be 1.0, got %f", analysis.channelWidth)
	// }
	// if analysis.channelLength != 6.0 {
	// 	t.Errorf("Expected channelLength to be 6.0, got %f", analysis.channelLength)
	// }
}

func TestNewPreCalculator(t *testing.T) {
	field := newMockField(6, 3)
	pc := NewPreCalculator(field)

	if pc.analysis == nil {
		t.Errorf("Expected analysis to be initialized")
	}
}
