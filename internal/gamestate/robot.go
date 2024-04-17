package gamestate

import (
	"fmt"

	"time"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"gonum.org/v1/gonum/mat"
)

type Team int

const (
	Blue Team = iota
	Yellow
)

type Robot struct {
	id          int
	team        Team
	pos         *mat.VecDense
	prevPos     *mat.VecDense
	lastUpdated time.Time
	vel         *mat.VecDense
	action      action.Action
}

func NewRobot(id int, team Team) *Robot {
	return &Robot{
		id:   id,
		team: team,
		pos:  mat.NewVecDense(3, []float64{0, 0, 0}), // in mm
		vel:  mat.NewVecDense(3, []float64{0, 0, 0}), // in mm/s
	}
}

func (r *Robot) SetPosition(x, y, w float64) {
	// Store the current time
	currentTime := time.Now()

	// Calculate the time difference in seconds
	timeDiff := currentTime.Sub(r.lastUpdated).Milliseconds()

	if timeDiff > 0 {
		// Calculate the change in position for x and y
		deltaX := x - r.pos.AtVec(0)
		deltaY := y - r.pos.AtVec(1)
		deltaW := w - r.pos.AtVec(2)

		// Calculate the velocity components
		velocityX := deltaX / (float64(timeDiff) * 1000)
		velocityY := deltaY / (float64(timeDiff) * 1000)
		velocityW := deltaW / (float64(timeDiff) * 1000)

		// Update the velocity using setVelocity
		r.setVelocity(velocityX, velocityY, velocityW)

		r.pos.SetVec(0, x)
		r.pos.SetVec(1, y)
		r.pos.SetVec(2, w)
	}

	// Update the lastUpdated time
	r.lastUpdated = currentTime
}

func (r *Robot) GetPosition() *mat.VecDense {
	return r.pos
}

func (r *Robot) setVelocity(v_x, v_y, v_w float64) {
	r.vel.SetVec(0, v_x)
	r.vel.SetVec(1, v_y)
	r.vel.SetVec(2, v_w)
}

func (r *Robot) String() string {
	x := r.pos.AtVec(0)
	y := r.pos.AtVec(1)
	w := r.pos.AtVec(2)

	v_x := r.pos.AtVec(0)
	v_y := r.pos.AtVec(1)
	v_w := r.pos.AtVec(2)

	posString := fmt.Sprintf("(%f, %f, %f)", x, y, w)
	velString := fmt.Sprintf("(%f, %f, %f)", v_x, v_y, v_w)

	return fmt.Sprintf("id: %d, pos: %s, vel: %s", r.id, posString, velString)
}

func (r *Robot) ToDTO() RobotDTO {
	return RobotDTO{
		Id:   r.id,
		Team: r.team,
		PosX: int(r.pos.AtVec(0)),
		PosY: int(r.pos.AtVec(1)),
		PosW: r.pos.AtVec(2),
		VelX: r.vel.AtVec(0),
		VelY: r.vel.AtVec(1),
		VelW: r.vel.AtVec(2),
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
