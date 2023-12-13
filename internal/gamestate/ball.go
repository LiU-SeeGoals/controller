package gamestate

import (
	"gonum.org/v1/gonum/mat"
)

type Ball struct {
	pos *mat.VecDense
	vel *mat.VecDense
}

func NewBall() *Ball {
	return &Ball{
		pos: mat.NewVecDense(3, []float64{0, 0, 0}),
		vel: mat.NewVecDense(3, []float64{0, 0, 0}),
	}
}

func (b *Ball) SetPosition(x, y, w float64) {
	b.pos.SetVec(0, x)
	b.pos.SetVec(1, y)
	b.pos.SetVec(2, w)
}

func (b *Ball) SetVelocity(v_x, v_y, v_w float64) {
	b.vel.SetVec(0, v_x)
	b.vel.SetVec(1, v_y)
	b.vel.SetVec(2, v_w)
}

func (b *Ball) ToDTO() BallDTO {
	return BallDTO{
		PosX: int(b.pos.AtVec(0)),
		PosY: int(b.pos.AtVec(1)),
		PosW: b.pos.AtVec(2),
		VelX: int(b.vel.AtVec(0)),
		VelY: int(b.vel.AtVec(1)),
		VelW: b.vel.AtVec(1),
	}
}

type BallDTO struct {
	PosX int
	PosY int
	PosW float64
	VelX int
	VelY int
	VelW float64
}