package info

import (
	"fmt"
	"time"

)

func NewBall(historyCapacity int) *Ball {
	return &Ball{
		history:         make([]BallState, historyCapacity),
		historyCapacity: historyCapacity,
	}
}

type BallPos struct {
	pos  Position
	time int64
}

type BallState struct {
	Valid      bool
	Velocity   Position
	Position   Position
	Visibility float64
	Timestamp  int64
	Source     string
}

func NewBallState(position, velocity Position, visibility float64, timestamp int64, source string) BallState {
	return BallState{
		Valid:      true,
		Position:   position,
		Velocity:   velocity,
		Visibility: visibility,
		Timestamp:  timestamp,
		Source:     source,
	}
}

func (b *BallState) String() string {
	return fmt.Sprintf("Position: %v, Velocity: %v, Visibility: %v, Timestamp: %v", b.Position, b.Velocity, b.Visibility, b.Timestamp)
}

type Ball struct {
	history         []BallState
	historyCapacity int
	writeIndex      int
}

// ############# Setters #############
func (b *Ball) SetState(state BallState) {
	b.history[b.writeIndex] = state
	b.writeIndex = (b.writeIndex + 1) % b.historyCapacity
}


// ############# Getters #############

func (b *Ball) getCurrentState() BallState {
	readIndex := (b.writeIndex - 1 + b.historyCapacity) % b.historyCapacity
	return b.history[readIndex]
}

func (b *Ball) GetPositionTime() (Position, int64, error) {
	currentState := b.getCurrentState()
	if !currentState.Valid {
		return Position{}, 0, fmt.Errorf("No position in history for ball")
	}

	return currentState.Position, currentState.Timestamp, nil
}


func (b *Ball) GetVisibility() float64 {
	state := b.getCurrentState()
	return state.Visibility
}

// get age
func (b *Ball) GetAge() int64 {
	_, ballTime, err := b.GetPositionTime()
	if err != nil {
		return 0
	}

	return time.Now().UnixMilli() - ballTime
}

func (b *Ball) GetPosition() (Position, error) {
	pos, _, err := b.GetPositionTime()

	return pos, err
}

func (b *Ball) GetVelocity() Position {
	state := b.getCurrentState()
	return state.Velocity
}

// ############# DTO #############

type BallDTO struct {
	PosX float64
	PosY float64
	PosZ float64
	VelX float64
	VelY float64
	VelZ float64
}

func (b *Ball) ToDTO() BallDTO {
	pos, _ := b.GetPosition()
	vel := b.GetVelocity()
	return BallDTO{
		PosX: pos.X,
		PosY: pos.Y,
		PosZ: pos.Z,
		VelX: vel.X,
		VelY: vel.Y,
		VelZ: vel.Z,
	}
}
