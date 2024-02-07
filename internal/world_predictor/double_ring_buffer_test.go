package world_predictor

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDoubleRingBufferCreation(t *testing.T) {
	drb := NewDoubleRingBuffer()
	assert.NotNil(t, drb, "Should not be nil")
	assert.Equal(t, 0, drb.GetLength(), "Should be empty initially")
	assert.NotNil(t, drb.GetGameStateInProgress(), "Game state in progress should not be nil")
}

func TestDoubleRingBufferPlaceGameState(t *testing.T) {
	drb := NewDoubleRingBuffer()
	for i := 0; i < bufferSize; i++ {
		drb.GetGameStateInProgress().SetFrameNumber(i)
		drb.PlaceGameState()
	}

	assert.Equal(t, 0, drb.GetLength(), "We have not update yet, so should be empty")
	drb.Update()
	assert.Equal(t, bufferSize, drb.GetLength(), "Should have buffer size amount of game states")

	for i := 0; i < bufferSize; i++ {
		assert.Equal(t, bufferSize - (i + 1), drb.GetGameState(i).GetFrameNumber(), "Should be in order")
	}
}

func TestDoubleRingBufferMultipleThreads(t *testing.T) {
	amountOfGameStates := 100
	count := 0
	drb := NewDoubleRingBuffer()
	go produceGameState(drb, amountOfGameStates, true)

	for count < amountOfGameStates {

		wait := rand.Intn(1) + 2
		time.Sleep(time.Duration(wait) * time.Nanosecond)

		drb.Update()
		 
		for i := drb.GetLength() - 1; i > -1; i-- {
			assert.Equal(t, count, drb.GetGameState(i).GetFrameNumber(), fmt.Sprintf("Should be in order, count: %d", count))
			count++
		}
	}
}

func TestDoubleRingBufferNewUpdateAtFront(t *testing.T) {
	amountOfGameStates := 100
	drb := NewDoubleRingBuffer()
	produceGameState(drb, amountOfGameStates, false)
	drb.Update()
	for i := 0; i < drb.GetLength(); i++ {
		assert.Equal(t, amountOfGameStates - i - 1, drb.GetGameState(i).GetFrameNumber(), fmt.Sprintf("Should be in order, count: %d", i))	
	}
}

func produceGameState(drb *DoubleRingBuffer, amount int, wait bool) {
	for j := 0; j < amount; j++ {
		drb.GetGameStateInProgress().SetFrameNumber(j)
		drb.PlaceGameState()
		if wait {
			wait := rand.Intn(1) + 1
			time.Sleep(time.Duration(wait) * time.Millisecond)
		}
	}
}