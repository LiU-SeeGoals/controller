package gamestate

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

type Team int

const (
	Blue Team = iota
	Yellow
)

type Robot struct {
	id   int
	team Team
	pos  *mat.VecDense
	vel  *mat.VecDense
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
	r.pos.SetVec(0, x)
	r.pos.SetVec(1, y)
	r.pos.SetVec(2, w)
}

func (r *Robot) SetVelocity(v_x, v_y, v_w float64) {
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

// GetID returns the id of the Robot.
func (r *Robot) GetID() int {
	return r.id
}

// GetTeam returns the team of the Robot.
func (r *Robot) GetTeam() Team {
	return r.team
}

// GetPosition returns the position vector of the Robot.
func (r *Robot) GetPosition() *mat.VecDense {
	return r.pos
}

// GetVelocity returns the velocity vector of the Robot.
func (r *Robot) GetVelocity() *mat.VecDense {
	return r.vel
}
