package info

import (
	"container/list"
)

type Ball struct {
	rawBall
	possessingRobot *Robot
}

func NewBall(historyCapacity int) *Ball {
	return &Ball{
		rawBall: rawBall{
			history:         list.New(),
			historyCapacity: historyCapacity,
		},
	}
}

func (b *Ball) GetVelocity() Position {

	if b.history.Len() < 2 {
		return Position{0, 0, 0, 0}
	}

	element := b.history.Front()
	ball := element.Value.(*rawBallPos)

	sum_deltas := Position{}

	for e := b.history.Front().Next(); e != nil; e = e.Next() {
		ball2 := e.Value.(*rawBallPos)
		dPos := ball2.pos.Sub(&ball.pos)
		dt := float64(ball2.time - ball.time)
		scaled := dPos.Scale(1 / dt)
		sum_deltas = sum_deltas.Add(&scaled)
	}
	return sum_deltas.Scale(1 / float64(b.history.Len()-1))

}

func (b *Ball) ToDTO() BallDTO {
	pos, _ := b.GetPosition()
	vel := b.GetVelocity()
	return BallDTO{
		PosX: pos.X,
		PosY: pos.Y,
		PosZ: pos.Z,
		VelX: vel.X,
		VelY: vel.Y,
		VelZ: vel.Z,
	}
}

type BallDTO struct {
	PosX float64
	PosY float64
	PosZ float64
	VelX float64
	VelY float64
	VelZ float64
}
