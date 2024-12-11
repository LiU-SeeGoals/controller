package state

import (
	"container/list"
)

type BallPos struct {
	pos  Position
	time int64
}

type Ball struct {
	history         *list.List
	historyCapacity int
}

func NewBall(historyCapacity int) *Ball {
	return &Ball{
		history:         list.New(),
		historyCapacity: historyCapacity,
	}
}

func (b *Ball) SetPositionTime(x, y, z float32, time int64) {
	if b.history.Len() >= b.historyCapacity {
		element := b.history.Back()
		b.history.Remove(element)

		ball := element.Value.(*BallPos)

		ball.pos.X = x
		ball.pos.Y = y
		ball.pos.Z = z
		ball.time = time

		b.history.PushFront(ball)
	} else {
		pos := Position{x, y, z, 0}
		b.history.PushFront(&BallPos{pos, time})
	}
}

func (b *Ball) GetPositionTime() (Position, int64) {
	if b.history.Len() == 0 {
		panic("No position in history for ball")
	}
	ball := b.history.Front().Value.(*BallPos)

	return ball.pos, ball.time
}

func (b *Ball) GetPosition() Position {
	pos, _ := b.GetPositionTime()
	return pos
}

func (b *Ball) GetVelocity() Position {

	if b.history.Len() < 2 {
		return Position{0, 0, 0, 0}
	}

	element := b.history.Front()
	ball := element.Value.(*BallPos)

	sum_deltas := Position{}

	for e := b.history.Front().Next(); e != nil; e = e.Next() {
		ball2 := e.Value.(*BallPos)
		dPos := ball2.pos.Sub(&ball.pos)
		dt := float32(ball2.time - ball.time)
		scaled := dPos.Scale(1 / dt)
		sum_deltas = sum_deltas.Add(&scaled)
	}
	return sum_deltas.Scale(1 / float32(b.history.Len()-1))

}

func (b *Ball) ToDTO() BallDTO {
	pos := b.GetPosition()
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
	PosX float32
	PosY float32
	PosZ float32
	VelX float32
	VelY float32
	VelZ float32
}
