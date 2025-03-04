package info

import (
	"container/list"
	"fmt"
	"math"
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

type rawRobotPos struct {
	pos  Position
	time int64
}

type Robot struct {
	rawRobot
}

type rawRobot struct {
	active          bool
	id              ID
	team            Team
	history         *list.List
	historyCapacity int
}

func NewRobot(id ID, team Team, history_capasity int) *Robot {
	return &Robot{
		rawRobot: rawRobot{
			active:          false,
			id:              id,
			team:            team,
			history:         list.New(),
			historyCapacity: history_capasity,
		},
	}
}

func (r *rawRobot) SetPositionTime(x, y, angle float32, time int64) {
	r.active = true
	if r.history.Len() >= r.historyCapacity {
		element := r.history.Back()
		r.history.Remove(element)

		robot := element.Value.(*rawRobotPos)

		robot.pos.X = x
		robot.pos.Y = y
		robot.pos.Angle = angle
		robot.time = time

		r.history.PushFront(robot)
	} else {
		pos := Position{x, y, 0, angle}
		r.history.PushFront(&rawRobotPos{pos, time})
	}
}

<<<<<<< HEAD
func (r *Robot) GetPositionTime() (Position, int64, error) {
=======
func (r *rawRobot) GetPositionTime() (Position, int64) {
>>>>>>> 4eb4ea5 (embeded rawRobot and rawBall in Robot and Ball)
	if r.history.Len() == 0 {
		return Position{}, 0, fmt.Errorf("No position in history for robot %d %s", r.id, r.team.String())
		// panic("No position in history for robot " + fmt.Sprint(r.id) + " " + r.team.String())
	}

	element := r.history.Front()
<<<<<<< HEAD
	robot := element.Value.(*RobotPos)
	return robot.pos, robot.time, nil
}

func (r *Robot) GetPosition() (Position, error) {
	pos, _, err := r.GetPositionTime()
	return pos, err
=======
	robot := element.Value.(*rawRobotPos)
	return robot.pos, robot.time
}

func (r *rawRobot) GetPosition() Position {
	pos, _ := r.GetPositionTime()
	return pos
>>>>>>> 4eb4ea5 (embeded rawRobot and rawBall in Robot and Ball)
}

func (r *rawRobot) FacingPosition(pos Position, threshold float64) bool {

	robotPos, err := r.GetPosition()
	if err != nil {
		return false
	}

	dx := pos.X - robotPos.X
	dy := pos.Y - robotPos.Y

	targetDirection := float32(math.Atan2(float64(dy), float64(dx)))
<<<<<<< HEAD
	currentDirection := robotPos.Angle
	
=======
	currentDirection := r.GetPosition().Angle

>>>>>>> 4eb4ea5 (embeded rawRobot and rawBall in Robot and Ball)
	angleDiff := math.Abs(float64(targetDirection - currentDirection))
	if angleDiff < threshold {
		return true
	} else {
		return false
	}
}
func (r *Robot) GetVelocity() Position {
	if r.history.Len() < 2 {
		return Position{0, 0, 0, 0}
	}

	element := r.history.Front()
	robot := element.Value.(*rawRobotPos)

	sum_deltas := Position{}

	for e := r.history.Front().Next(); e != nil; e = e.Next() {
		robot2 := e.Value.(*rawRobotPos)
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

		robot1 := f.Value.(*rawRobotPos)
		robot2 := s.Value.(*rawRobotPos)
		robot3 := t.Value.(*rawRobotPos)

		vel1 := robot2.pos.Sub(&robot1.pos)
		vel2 := robot3.pos.Sub(&robot2.pos)

		dist1 := vel1.Norm()
		dist2 := vel2.Norm()

		dt := float32(robot3.time - robot1.time)

		accelerations += (dist2 - dist1) / dt

	}

	return accelerations / float32(r.history.Len()-1)
}

func (r *rawRobot) SetActive(active bool) {
	r.active = active
}

func (r *rawRobot) IsActive() bool {
	return r.active
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

func (r *rawRobot) GetID() ID {
	return r.id
}

type RobotDTO struct {
	Id       ID
	Team     Team
	X        float32
	Y        float32
	Angle    float32
	VelX     float32
	VelY     float32
	VelAngle float32
}
