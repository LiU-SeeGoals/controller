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

// dont use outside of gamestate/worldpredictor
func (b *Ball) GetPosition() *mat.VecDense {
	return b.pos
}

// dont use outside of gamestate/worldpredictor
func (b *Ball) SetPosition(x, y, w float64) {
	b.pos.SetVec(0, x)
	b.pos.SetVec(1, y)
	b.pos.SetVec(2, w)
}

// dont use outside of gamestate/worldpredictor
func (b *Ball) SetVelocity(v_x, v_y, v_w float64) {
	b.vel.SetVec(0, v_x)
	b.vel.SetVec(1, v_y)
	b.vel.SetVec(2, v_w)
}

func (b *Ball) NormalizePosition(normalizationFactor float64) {
	b.pos.SetVec(0, b.pos.AtVec(0)/normalizationFactor)
	b.pos.SetVec(1, b.pos.AtVec(1)/normalizationFactor)
	b.pos.SetVec(2, b.pos.AtVec(2)/normalizationFactor)
}