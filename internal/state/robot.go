package state

import (
	"container/list"
	"fmt"
)

type Team int8
type ID uint8
type RobotTeam [TEAM_SIZE]*Robot

const (
	Blue Team = iota
	Yellow
)

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
		id:              id,
		team:            team,
		history:         list.New(),
		historyCapacity: history_capasity,
	}
}

func (r *Robot) SetPositionTime(x, y, angel float32, time int64) {
	if r.history.Len() >= r.historyCapacity {
		element := r.history.Back()
		r.history.Remove(element)

		robot := element.Value.(*RobotPos)

		robot.pos.X = x
		robot.pos.Y = y
		robot.pos.Angel = angel
		robot.time = time

		r.history.PushFront(robot)
	} else {
		pos := Position{x, y, 0, angel}
		r.history.PushFront(&RobotPos{pos, time})
	}
}

func (r *Robot) GetPositionTime() (Position, int64) {
	if r.history.Len() == 0 {
		panic("No position in history")
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
		dPos := robot2.pos.Sub(robot.pos)
		dt := float32(robot2.time - robot.time)
		// TODO: lets add exponential decay so that the most recent deltas have more weight
		sum_deltas = sum_deltas.Add(dPos.Scale(1 / dt))
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

		vel1 := robot2.pos.Sub(robot1.pos)
		vel2 := robot3.pos.Sub(robot2.pos)

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

	posString := fmt.Sprintf("(%f, %f, %f)", pos.X, pos.Y, pos.Angel)
	velString := fmt.Sprintf("(%f, %f, %f)", vel.X, vel.Y, vel.Angel)

	return fmt.Sprintf("id: %d, pos: %s, vel: %s", r.id, posString, velString)
}

func (r *Robot) ToDTO() RobotDTO {
	pos := r.GetPosition()
	vel := r.GetVelocity()

	return RobotDTO{
		Id:       r.id,
		Team:     r.team,
		X:        pos.X,
		Y:        pos.Y,
		Angel:    pos.Angel,
		VelX:     vel.X,
		VelY:     vel.Y,
		VelAngel: vel.Angel,
	}
}

func (r *Robot) GetID() ID {
	return r.id
}

type RobotDTO struct {
	Id       ID
	Team     Team
	X        float32
	Y        float32
	Angel    float32
	VelX     float32
	VelY     float32
	VelAngel float32
}
