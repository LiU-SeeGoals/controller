package world_predictor

import (
	"sync"
	"time"

	"github.com/LiU-SeeGoals/controller/internal/gamestate"
)

type DoubleRingBuffer struct {
	buffers [2]*RingBuffer
	active int32
	gameStateInProgress *gamestate.GameState
	lock sync.Mutex
}

func NewDoubleRingBuffer() *DoubleRingBuffer {
	return &DoubleRingBuffer{
		buffers: [2]*RingBuffer{
			NewRingBuffer(),
			NewRingBuffer(),
		},
		active: 0,
		gameStateInProgress: gamestate.NewGameState(),
		lock: sync.Mutex{},
	}
}

func (drb *DoubleRingBuffer) Update() {
	for drb.getInactiveBuffer().GetLength() <= 0 {
		// Wait for the inactive buffer to be filled
		// Using the back off strategy
		time.Sleep(time.Millisecond)
	}

	drb.lock.Lock()
	drb.active = (drb.active + 1) % 2
	drb.getInactiveBuffer().Clear()
	drb.lock.Unlock()
}

func (drb *DoubleRingBuffer) PlaceGameState() {
	drb.lock.Lock()
	drb.gameStateInProgress = drb.getInactiveBuffer().placeNewGameState(drb.gameStateInProgress)
	drb.lock.Unlock()
	drb.gameStateInProgress.Clear()
}

func (drb *DoubleRingBuffer) GetGameStateInProgress() *gamestate.GameState {
	return drb.gameStateInProgress
}

func (drb *DoubleRingBuffer) GetGameState(index int) *gamestate.GameState {
	return drb.buffers[drb.active].GetGameState(index)
}

func (drb *DoubleRingBuffer) GetLength() int {
	return drb.buffers[drb.active].GetLength()
}

func (drb *DoubleRingBuffer) getActiveBuffer() *RingBuffer {
	return drb.buffers[drb.active]
}

func (drb *DoubleRingBuffer) getInactiveBuffer() *RingBuffer {
	return drb.buffers[(drb.active + 1) % 2]
}