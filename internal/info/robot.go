package info

import (
	"fmt"
	"math"
	"time"

	. "github.com/LiU-SeeGoals/controller/internal/logger"
)

type ID uint8

const historyCapasity int = 10

func NewRobot(id ID, team Team, historyCapacity int) *Robot {
	return &Robot{
		active:          false,
		id:              id,
		team:            team,
		history:	 make([]RobotState, historyCapacity),
		historyCapacity: historyCapacity,
	}
}

type RobotPos struct {
	pos  Position
	time int64
}

type RobotState struct {
	Valid      bool
	velocity   Position
	angularVel float64
	position   Position
	visibility float64
	timestamp  int64
	source     string
}

func NewRobotState(position, velocity Position, visibility float64, timestamp int64, source string) RobotState {
	return RobotState{
		Valid:      true,
		position:   position,
		velocity:   velocity,
		visibility: visibility,
		timestamp:  timestamp,
		source:     source,
	}
}

func (r *RobotState) String() string {
	return fmt.Sprintf("Position: %v, Velocity: %v, Angular Velocity: %v, Visibility: %v, Timestamp: %v", r.position, r.velocity, r.angularVel, r.visibility, r.timestamp)
}


type Robot struct {
	// RobotState
	active          bool
	id              ID
	team            Team
	history         []RobotState
	writeIndex      int
	historyCapacity int
}

func (r *Robot) SetState(state RobotState) {
	r.history[r.writeIndex] = state
	r.writeIndex = (r.writeIndex + 1) % r.historyCapacity
	// if r.history.Len() >= r.historyCapacity {
	// 	element := r.history.Back()
	// 	r.history.Remove(element)
	//
	// 	robotState := element.Value.(*RobotState)
	// 	robotState = &state
	//
	// 	r.history.PushFront(robotState)
	// 	fmt.Println("Robot history is full")
	// } else {
	// 	fmt.Println("Robot history is not full")
	// 	r.history.PushFront(&state)
	// }
}

// // SetVelocity sets the robot's velocity.
// func (r *Robot) SetVelocity(x, y, angle float64) {
// 	r.velocity = Position{x, y, 0, angle}
// }
//
// // SetVisibility sets the robot's visibility.
// func (r *Robot) SetVisibility(visibility float64) {
// 	r.visibility = visibility
// }

// func (r *Robot) SetPositionTime(x, y, angle float64, time int64) {
// 	r.active = true
// 	if r.history.Len() >= r.historyCapacity {
// 		element := r.history.Back()
// 		r.history.Remove(element)
//
// 		robot := element.Value.(*RobotPos)
//
// 		robot.pos.X = x
// 		robot.pos.Y = y
// 		robot.pos.Angle = angle
// 		robot.time = time
//
// 		r.history.PushFront(robot)
// 	} else {
// 		pos := Position{x, y, 0, angle}
// 		r.history.PushFront(&RobotPos{pos, time})
// 	}
// }

func (r *Robot) GetPositionTime() (Position, int64, error) {
	currentState := r.getCurrentState()
	if !currentState.Valid {
		return Position{}, 0, fmt.Errorf("No position in history for robot %d %s", r.id, r.team.String())
		// panic("No position in history for robot " + fmt.Sprint(r.id) + " " + r.team.String())
	}

	return currentState.position, currentState.timestamp, nil
}

func (r *Robot) GetPosition() (Position, error) {
	pos, _, err := r.GetPositionTime()
	return pos, err
}

func (r *Robot) GetID() ID {
	return r.id
}

func (r *Robot) getCurrentState() RobotState {
	readIndex := (r.writeIndex - 1 + r.historyCapacity) % r.historyCapacity
	return r.history[readIndex]
}

func (r *Robot) IsActive() bool {
	age := time.Now().UnixMilli() - r.getCurrentState().timestamp

	if age > 10000 {
		return false
	} else {
		return true
	}
}

func (r *Robot) At(pos Position, threshold float64) bool {
	robotPos, err := r.GetPosition()
	if err != nil {
		Logger.Errorf("Position retrieval failed - Robot: %v\n", err)
		return false
	}

	return robotPos.Distance(pos) < threshold
}

func (r *Robot) DribblerPos() Position {

	robotPos, _ := r.GetPosition()
	robotPos.X += 90 * math.Cos(robotPos.Angle) // WARN: Magic number
	robotPos.Y += 90 * math.Sin(robotPos.Angle) // WARN: Magic number
	return robotPos
}

func (r *Robot) Facing(target Position, threshold float64) bool {
	pos, err := r.GetPosition()
	if err != nil {
		return false
	}
	targetDirection := pos.AngleToPosition(target)
	currentDirection := pos.Angle

	angleDiff := math.Abs(float64(targetDirection - currentDirection))
	if angleDiff < threshold {
		return true
	} else {
		return false
	}

}

func (r *Robot) GetVelocity() Position {
	state := r.getCurrentState()
	if !state.Valid {
		return Position{}
		// panic("No position in history for robot " + fmt.Sprint(r.id) + " " + r.team.String())
	}

	return state.velocity
}

// func (r *Robot) GetAcceleration() float64 {
// 	if r.history.Len() < 3 {
// 		return float64(0) // Not enough data points to calculate acceleration
// 	}
//
// 	accelerations := float64(0)
// 	for f, s, t := r.history.Front(), r.history.Front().Next(), r.history.Front().Next().Next(); t != nil; f, s, t = f.Next(), s.Next(), t.Next() {
//
// 		robot1 := f.Value.(*RobotPos)
// 		robot2 := s.Value.(*RobotPos)
// 		robot3 := t.Value.(*RobotPos)
//
// 		vel1 := robot2.pos.Sub(&robot1.pos)
// 		vel2 := robot3.pos.Sub(&robot2.pos)
//
// 		dist1 := vel1.Length()
// 		dist2 := vel2.Length()
//
// 		dt := float64(robot3.time - robot1.time)
//
// 		accelerations += (dist2 - dist1) / dt
//
// 	}
//
// 	return accelerations / float64(r.history.Len()-1)
// }

func (r *Robot) GetTeam() Team {
	return r.team
}

func (r *Robot) String() string {

	pos, err := r.GetPosition()

	if err != nil {
		return ""
	}

	vel := r.GetVelocity()

	posString := fmt.Sprintf("(%f, %f, %f)", pos.X, pos.Y, pos.Angle)
	velString := fmt.Sprintf("(%f, %f, %f)", vel.X, vel.Y, vel.Angle)

	return fmt.Sprintf("id: %d, pos: %s, vel: %s", r.id, posString, velString)
}

func (r *Robot) ToDTO() RobotDTO {
	if !r.active {
		return RobotDTO{}
	}
	pos, _ := r.GetPosition()
	vel := r.GetVelocity()

	return RobotDTO{
		Id:       r.id,
		Team:     r.team,
		X:        pos.X,
		Y:        pos.Y,
		Angle:    pos.Angle,
		VelX:     vel.X,
		VelY:     vel.Y,
		VelAngle: vel.Angle,
	}
}

type RobotDTO struct {
	Id       ID
	Team     Team
	X        float64
	Y        float64
	Angle    float64
	VelX     float64
	VelY     float64
	VelAngle float64
}
