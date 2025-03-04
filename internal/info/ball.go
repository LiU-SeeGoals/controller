package info

import (
	"container/list"
	"errors"
)

type Ball struct {
	rawBall
	possessingRobot *Robot
}

type rawBallPos struct {
	pos  Position
	time int64
}

type rawBall struct {
	history         *list.List
	historyCapacity int
}

func NewBall(historyCapacity int) *Ball {
	return &Ball{
		rawBall: rawBall{
			history:         list.New(),
			historyCapacity: historyCapacity,
		},
	}
}

func (b *rawBall) SetPositionTime(x, y, z float32, time int64) {
	if b.history.Len() >= b.historyCapacity {
		element := b.history.Back()
		b.history.Remove(element)

		ball := element.Value.(*rawBallPos)

		ball.pos.X = x
		ball.pos.Y = y
		ball.pos.Z = z
		ball.time = time

		b.history.PushFront(ball)
	} else {
		pos := Position{x, y, z, 0}
		b.history.PushFront(&rawBallPos{pos, time})
	}
}

<<<<<<< HEAD
func (b *Ball) GetPositionTime() (Position, int64, error) {
=======
func (b *rawBall) GetPositionTime() (Position, int64) {
>>>>>>> 4eb4ea5 (embeded rawRobot and rawBall in Robot and Ball)
	if b.history.Len() == 0 {
		return Position{}, 0, errors.New("No position in history for ball")
	}
	ball := b.history.Front().Value.(*rawBallPos)

	return ball.pos, ball.time, nil
}

<<<<<<< HEAD
func (b *Ball) GetPosition() (Position, error) {
	pos, _, err := b.GetPositionTime()
	return pos, err
=======
func (b *rawBall) GetPosition() Position {
	pos, _ := b.GetPositionTime()
	return pos
>>>>>>> 4eb4ea5 (embeded rawRobot and rawBall in Robot and Ball)
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
		dt := float32(ball2.time - ball.time)
		scaled := dPos.Scale(1 / dt)
		sum_deltas = sum_deltas.Add(&scaled)
	}
	return sum_deltas.Scale(1 / float32(b.history.Len()-1))

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
	PosX float32
	PosY float32
	PosZ float32
	VelX float32
	VelY float32
	VelZ float32
}
