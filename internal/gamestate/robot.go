package gamestate

import (
	"container/list"
	"fmt"

	"gonum.org/v1/gonum/mat"
)

type Team int

const (
	Blue Team = iota
	Yellow
)

type RobotPos struct {
	pos  *mat.VecDense
	time int64
}

type Robot struct {
	id               int
	team             Team
	history          *list.List
	history_capacity int
	vel              *mat.VecDense // in mm/s
	max_speed        float64       // in mm/s
}

func NewRobot(id int, team Team, history_capasity int) *Robot {
	return &Robot{
		id:               id,
		team:             team,
		history:          list.New(),
		history_capacity: history_capasity,
		vel:              mat.NewVecDense(3, []float64{0, 0, 0}), // in mm/s
		max_speed:        0,
	}
}

func (r *Robot) SetPositionTime(x, y, w float64, time int64) {
	if r.history.Len() >= r.history_capacity {
		element := r.history.Back()
		r.history.Remove(element)

		robot := element.Value.(*RobotPos)

		robot.pos.SetVec(0, x)
		robot.pos.SetVec(1, y)
		robot.pos.SetVec(2, w)
		robot.time = time

		r.history.PushFront(robot)
	} else {
		pos := mat.NewVecDense(3, []float64{x, y, w})
		r.history.PushFront(&RobotPos{pos, time})
	}
}

func (r *Robot) GetPositionTime() (*mat.VecDense, int64) {
	if r.history.Len() == 0 {
		return nil, 0
	}

	element := r.history.Front()
	robot := element.Value.(*RobotPos)
	return robot.pos, robot.time
}

func (r *Robot) GetPosition() *mat.VecDense {
	pos, _ := r.GetPositionTime()
	return pos
}

func (r *Robot) UpdateVelocity() {
	if r.history.Len() < 2 {
		return
	}

	robot1 := r.history.Front().Value.(*RobotPos)
	robot2 := r.history.Front().Next().Value.(*RobotPos)

	dt := float64(robot2.time - robot1.time)
	if dt > 0 {
		dx := robot2.pos.AtVec(0) - robot1.pos.AtVec(0)
		dy := robot2.pos.AtVec(1) - robot1.pos.AtVec(1)
		dw := robot2.pos.AtVec(2) - robot1.pos.AtVec(2)

		v_x := dx / dt
		v_y := dy / dt
		v_w := dw / dt

		r.vel.SetVec(0, v_x)
		r.vel.SetVec(1, v_y)
		r.vel.SetVec(2, v_w)
	} else {
		r.vel.SetVec(0, 0)
		r.vel.SetVec(1, 0)
		r.vel.SetVec(2, 0)
	}
}

func (r *Robot) GetVelocity() *mat.VecDense {
	return r.vel
}

func (r *Robot) String() string {

	pos := r.GetPosition()
	x := pos.AtVec(0)
	y := pos.AtVec(1)
	w := pos.AtVec(2)

	vel := r.GetVelocity()
	v_x := vel.AtVec(0)
	v_y := vel.AtVec(1)
	v_w := vel.AtVec(2)

	posString := fmt.Sprintf("(%f, %f, %f)", x, y, w)
	velString := fmt.Sprintf("(%f, %f, %f)", v_x, v_y, v_w)

	return fmt.Sprintf("id: %d, pos: %s, vel: %s", r.id, posString, velString)
}

func (r *Robot) ToDTO() RobotDTO {
	pos := r.GetPosition()
	if pos == nil {
		return RobotDTO{
			Id:   r.id,
			Team: r.team,
			PosX: 0,
			PosY: 0,
			PosW: 0,
			VelX: 0,
			VelY: 0,
			VelW: 0,
		}

	}

	vel := r.GetVelocity()

	return RobotDTO{
		Id:   r.id,
		Team: r.team,
		PosX: int(pos.AtVec(0)),
		PosY: int(pos.AtVec(1)),
		PosW: pos.AtVec(2),
		VelX: vel.AtVec(0),
		VelY: vel.AtVec(1),
		VelW: vel.AtVec(2),
	}
}

func (r *Robot) GetID() int {
	return r.id
}

type RobotDTO struct {
	Id   int
	Team Team
	PosX int
	PosY int
	PosW float64
	VelX float64
	VelY float64
	VelW float64
}
