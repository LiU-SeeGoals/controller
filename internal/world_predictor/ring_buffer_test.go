package world_predictor

import (
	"testing"

	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"github.com/stretchr/testify/assert" // Optional, for more expressive assertions
)

func TestRingBuffer(t *testing.T) {
    rb := NewRingBuffer()

    // Test initial state
    assert.Equal(t, 0, rb.GetLength(), "New ring buffer should be empty")

	// Test passing placing
	newGameState := gamestate.NewGameState()
	newGameState.SetFrameNumber(0)
	workingGameState := rb.placeNewGameState(newGameState)
	assert.Equal(t, newGameState, rb.GetGameState(0), "Should return the gamestate we just placed.")

	// Test clearing
	rb.Clear()
	assert.Equal(t, 0, rb.GetLength(), "Ring buffer should be empty after clearing")


    // Test adding game states
    for i := 0; i < bufferSize; i++ {
		workingGameState.SetFrameNumber(i)
        workingGameState = rb.placeNewGameState(workingGameState)
        assert.Equal(t, i+1, rb.GetLength(), "Ring buffer length should increase")
    }

	// Checking the order
	for i := 0; i < bufferSize; i++ {
		assert.Equal(t, bufferSize - (i + 1), rb.GetGameState(i).GetFrameNumber(), "Ring buffer should be in order")
	}

	// Test adding game states after buffer is full
	for i := bufferSize; i < bufferSize*2; i++ {
		workingGameState.SetFrameNumber(i)
		workingGameState = rb.placeNewGameState(workingGameState)
		assert.Equal(t, bufferSize, rb.GetLength(), "Ring buffer length should not increase")
	}

	// Checking the order
	for i := 0; i < bufferSize; i++ {
		assert.Equal(t, bufferSize*2 - (i + 1), rb.GetGameState(i).GetFrameNumber(), "Ring buffer should be in order")
	}
}