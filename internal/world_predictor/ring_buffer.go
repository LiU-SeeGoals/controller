package world_predictor

import (
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
)

const bufferSize = 20

type RingBuffer struct {
	buffer []*gamestate.GameState
	index int
	amount int
}

func NewRingBuffer() *RingBuffer {
	buffer := make([]*gamestate.GameState, bufferSize)
	for i := 0; i < bufferSize; i++ {
		buffer[i] = gamestate.NewGameState()
	}

	return &RingBuffer{
		buffer: buffer,
		index: 0,
		amount: 0,
	}
}

func (rb *RingBuffer) placeNewGameState(newGameState *gamestate.GameState) *gamestate.GameState {
	nextIndex := rb.getNextIndex()
	previousGameState := rb.buffer[nextIndex]
	rb.buffer[nextIndex] = newGameState
	rb.index = rb.getNextIndex()
	rb.incrementAmount()
	return previousGameState
}

func (rb *RingBuffer) GetGameState(index int) *gamestate.GameState {
	return rb.buffer[rb.getIndex(index)]
}

func (rb *RingBuffer) GetLength() int {
	return rb.amount
}

func (rb *RingBuffer) Clear() {
	rb.amount = 0
}

func (rb *RingBuffer) getIndex(index int) int {
	return (rb.index - index + bufferSize) % bufferSize
}

func (rb *RingBuffer) getNextIndex() int {
	return (rb.index + 1 + bufferSize) % bufferSize
}

func (rb *RingBuffer) incrementAmount() {
	if rb.amount < bufferSize {
		rb.amount++
	}
}




