package state

import (
	"container/list"

	"gonum.org/v1/gonum/mat"
)

type BallPos struct {
	pos  *mat.VecDense
	time int64
}

type Ball struct {
	history         *list.List
	historyCapacity int
	vel             *mat.VecDense
	maxSpeed        float64
}

func NewBall(historyCapacity int) *Ball {
	return &Ball{
		history:         list.New(),
		historyCapacity: historyCapacity,
		vel:             mat.NewVecDense(3, []float64{0, 0, 0}),
		maxSpeed:        1,
	}
}

func (b *Ball) SetPositionTime(x, y, z float64, time int64) {
	if b.history.Len() >= b.historyCapacity {
		element := b.history.Back()
		b.history.Remove(element)

		ball := element.Value.(*BallPos)

		ball.pos.SetVec(0, x)
		ball.pos.SetVec(1, y)
		ball.pos.SetVec(2, z)
		ball.time = time

		b.history.PushFront(ball)
	} else {
		pos := mat.NewVecDense(3, []float64{x, y, z})
		b.history.PushFront(&BallPos{pos, time})
	}
}

func (b *Ball) GetPositionTime() (*mat.VecDense, int64) {
	ball := b.history.Front().Value.(*BallPos)

	return ball.pos, ball.time
}

func (b *Ball) GetPosition() *mat.VecDense {
	pos, _ := b.GetPositionTime()
	return pos
}

func (b *Ball) UpdateVelocity() {
	if b.history.Len() < 2 {
		return
	}

	ball1 := b.history.Front().Value.(*BallPos)
	ball2 := b.history.Front().Next().Value.(*BallPos)

	dt := float64(ball2.time - ball1.time)
	if dt > 0 {
		dx := ball2.pos.AtVec(0) - ball1.pos.AtVec(0)
		dy := ball2.pos.AtVec(1) - ball1.pos.AtVec(1)
		dz := ball2.pos.AtVec(2) - ball1.pos.AtVec(2)

		vX := dx / dt
		vY := dy / dt
		vZ := dz / dt

		b.vel.SetVec(0, vX)
		b.vel.SetVec(1, vY)
		b.vel.SetVec(2, vZ)

		vec := mat.NewVecDense(3, []float64{vX, vY, 0})
		speed := mat.Norm(vec, 2)
		if speed > b.maxSpeed {
			b.maxSpeed = speed
		}
	}
}

func (b *Ball) UpdateMaxSpeed() {
	speed := mat.Norm(b.vel, 2)
	if speed > b.maxSpeed {
		b.maxSpeed = speed
	}
}

func (b *Ball) GetVelocity() *mat.VecDense {
	return b.vel
}

func (b *Ball) ToDTO() BallDTO {
	pos := b.GetPosition()
	vel := b.GetVelocity()
	return BallDTO{
		PosX: int(pos.AtVec(0)),
		PosY: int(pos.AtVec(1)),
		PosW: pos.AtVec(2),
		VelX: int(vel.AtVec(0)),
		VelY: int(vel.AtVec(1)),
		VelW: vel.AtVec(1),
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
