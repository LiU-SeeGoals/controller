package gamestate

import (
	"github.com/LiU-SeeGoals/controller/internal/parsed_vision"
	"github.com/LiU-SeeGoals/proto_go/robot_action"
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

func (b *Ball) GetX() float64 {
	return b.pos.AtVec(0)
}

func (b *Ball) GetY() float64 {
	return b.pos.AtVec(1)
}

func (b *Ball) GetZ() float64 {
	return b.pos.AtVec(2)
}

func (b *Ball) GetVelocity() *mat.VecDense {
	return b.vel
}

func (b *Ball) GetVelX() float64 {
	return b.vel.AtVec(0)
}

func (b *Ball) GetVelY() float64 {
	return b.vel.AtVec(1)
}

func (b *Ball) GetVelZ() float64 {
	return b.vel.AtVec(2)
}

// dont use outside of gamestate/worldpredictor
func (b *Ball) SetPosition(x, y, z float64) {
	b.pos.SetVec(0, x)
	b.pos.SetVec(1, y)
	b.pos.SetVec(2, z)
}

// dont use outside of gamestate/worldpredictor
func (b *Ball) SetVelocity(v_x, v_y, v_z float64) {
	b.vel.SetVec(0, v_x)
	b.vel.SetVec(1, v_y)
	b.vel.SetVec(2, v_z)
}

func (b *Ball) GetParsedBall() *parsed_vision.Ball {
	return &parsed_vision.Ball{
		Pos: &robot_action.Vector2D{
			X: int32(b.GetX()),
			Y: int32(b.GetY()),
		},
		Vel: &robot_action.Vector2D{
			X: int32(b.GetVelX()),
			Y: int32(b.GetVelY()),
		},
	}

}

func (b *Ball) NormalizePosition(normalizationFactor float64) {
	b.pos.SetVec(0, b.pos.AtVec(0)/normalizationFactor)
	b.pos.SetVec(1, b.pos.AtVec(1)/normalizationFactor)
	b.pos.SetVec(2, b.pos.AtVec(2)/normalizationFactor)
}