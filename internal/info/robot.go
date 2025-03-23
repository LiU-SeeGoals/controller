package info

import (
	"container/list"
	"fmt"
)

type Team int8
type ID uint8
type RobotTeam [TEAM_SIZE]*Robot

const (
	UNKNOWN Team = 0
	Yellow  Team = 1
	Blue    Team = 2
)

func (t Team) String() string {
	switch t {
	case Yellow:
		return "Yellow"
	case Blue:
		return "Blue"
	default:
		return "UNKNOWN"
	}
}

type RobotPos struct {
	pos  Position
	time int64
}

type Robot struct {
	active          bool
	id              ID
	team            Team
	history         *list.List
	historyCapacity int
}

func NewRobot(id ID, team Team, history_capasity int) *Robot {
	return &Robot{
		active:          false,
		id:              id,
		team:            team,
		history:         list.New(),
		historyCapacity: history_capasity,
	}
}

func (r *Robot) SetPositionTime(x, y, angle float32, time int64) {
	r.active = true
	if r.history.Len() >= r.historyCapacity {
		element := r.history.Back()
		r.history.Remove(element)

		robot := element.Value.(*RobotPos)

		robot.pos.X = x
		robot.pos.Y = y
		robot.pos.Angle = angle
		robot.time = time

		r.history.PushFront(robot)
	} else {
		pos := Position{x, y, 0, angle}
		r.history.PushFront(&RobotPos{pos, time})
	}
}

func (r *Robot) GetPositionTime() (Position, int64) {
	if r.history.Len() == 0 {
		panic("No position in history for robot " + fmt.Sprint(r.id) + " " + r.team.String())
	}

	element := r.history.Front()
	robot := element.Value.(*RobotPos)
	return robot.pos, robot.time
}

func (r *Robot) GetPosition() Position {
	pos, _ := r.GetPositionTime()
	return pos
}

func (r *Robot) GetVelocity() Position {
	if r.history.Len() < 2 {
		return Position{0, 0, 0, 0}
	}

	element := r.history.Front()
	robot := element.Value.(*RobotPos)

	sum_deltas := Position{}

	for e := r.history.Front().Next(); e != nil; e = e.Next() {
		robot2 := e.Value.(*RobotPos)
		dPos := robot2.pos.Sub(&robot.pos)
		dt := float32(robot2.time - robot.time)
		// TODO: lets add exponential decay so that the most recent deltas have more weight
		scaled := dPos.Scale(1 / dt)
		sum_deltas = sum_deltas.Add(&scaled)
	}
	return sum_deltas.Scale(1 / float32(r.history.Len()-1))
}

func (r *Robot) GetAcceleration() float32 {
	if r.history.Len() < 3 {
		return float32(0) // Not enough data points to calculate acceleration
	}

	accelerations := float32(0)
	for f, s, t := r.history.Front(), r.history.Front().Next(), r.history.Front().Next().Next(); t != nil; f, s, t = f.Next(), s.Next(), t.Next() {

		robot1 := f.Value.(*RobotPos)
		robot2 := s.Value.(*RobotPos)
		robot3 := t.Value.(*RobotPos)

		vel1 := robot2.pos.Sub(&robot1.pos)
		vel2 := robot3.pos.Sub(&robot2.pos)

		dist1 := vel1.Norm()
		dist2 := vel2.Norm()

		dt := float32(robot3.time - robot1.time)

		accelerations += (dist2 - dist1) / dt

	}

	return accelerations / float32(r.history.Len()-1)
}

func (r *Robot) SetActive(active bool) {
	r.active = active
}

func (r *Robot) IsActive() bool {
	return r.active
}

func (r *Robot) String() string {

	pos := r.GetPosition()
	vel := r.GetVelocity()

	posString := fmt.Sprintf("(%f, %f, %f)", pos.X, pos.Y, pos.Angle)
	velString := fmt.Sprintf("(%f, %f, %f)", vel.X, vel.Y, vel.Angle)

	return fmt.Sprintf("id: %d, pos: %s, vel: %s", r.id, posString, velString)
}

func (r *Robot) ToDTO() RobotDTO {
	if !r.active {
		return RobotDTO{}
	}
	vel := r.GetVelocity()

	return RobotDTO{
		Id:       r.id,
		Team:     r.team,
		VelX:     vel.X,
		VelY:     vel.Y,
		VelAngle: vel.Angle,
	}
}

func (r *Robot) GetID() ID {
	return r.id
}

type RobotDTO struct {
	Id       ID
	Team     Team
	VelX     float32
	VelY     float32
	VelAngle float32
}
