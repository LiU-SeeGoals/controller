package info

import (
	"container/list"
	"errors"
)

type rawBallPos struct {
	pos  Position
	time int64
}

type rawBall struct {
	history         *list.List
	historyCapacity int
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

func (b *rawBall) GetPositionTime() (Position, int64, error) {
	if b.history.Len() == 0 {
		return Position{}, 0, errors.New("No position in history for ball")
	}
	ball := b.history.Front().Value.(*rawBallPos)

	return ball.pos, ball.time, nil
}

func (b *rawBall) GetPosition() (Position, error) {
	pos, _, err := b.GetPositionTime()

	return pos, err
}
