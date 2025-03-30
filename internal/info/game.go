package info

import (
	"math"
	"fmt"

	"github.com/LiU-SeeGoals/proto_go/ssl_vision"
)

// Maybe State should also be private, so we can keep track of coordinate system here?
type GameInfo struct {
	State  *GameState
	Status *GameStatus
	field  *ssl_vision.SSL_GeometryFieldSize
}

func NewGameInfo(capacity int) *GameInfo {
	return &GameInfo{
		State:  NewGameState(capacity),
		Status: NewGameStatus(),
		field:  nil,
	}
}

// Rotates all positions so that no matter 
// which side we are on we can use the same coordinate system
func correctedPosition(team Team, pos Position) Position{

	if team == Yellow{
		return pos
	} else if team == Blue{
		return pos.Rotate(math.Pi)
	} else {
		panic(fmt.Sprintf("Incorrect team given %v", team))
	}
}

func (gi GameInfo) HomeGoalLine(team Team) Position {

	x := float64(gi.field.GetFieldLength()/2 - gi.field.GetGoalWidth())

	pos := Position{X: x, Y: 0, Z: 0, Angle: 0}

	return correctedPosition(team, pos)
}

func (gi GameInfo) HasField() bool {

	if gi.field == nil{
		return false
	}

	return true
}
func (gi *GameInfo) SetField(field *ssl_vision.SSL_GeometryFieldSize) {
	gi.field = field
}

func (gi GameInfo) GoalArea(team Team) Position {

	x := float64(gi.field.GetFieldLength()/2 - gi.field.GetGoalWidth())

	pos := Position{X: x, Y: 0, Z: 0, Angle: 0}

	return correctedPosition(team, pos)
}
