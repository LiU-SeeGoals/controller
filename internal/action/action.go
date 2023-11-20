package action

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

type RealAction int

const (
	RealStop RealAction = 0
	RealKick = 1
	RealMove = 2
	RealInit = 3
)

type Action interface {
	TranslateReal() []byte
	TranslateGrsim() int
}

type Stop struct {
	Id int
}

type Move struct {
	Id   int
	Pos  *mat.VecDense
	Goal *mat.VecDense
}

type Kick struct {
	Id int
	Speed int
}

type Init struct {
	Id int
}

func (m *Move) TranslateReal() []byte {
	message := make([]byte, 19)
	message[0] = byte(19)
	message[1] = byte(RealMove)
	message[2] = byte(m.Id)

	appendPositionData(message, 3, m.Pos)
	appendPositionData(message, 11, m.Goal)

	return message
}

func (m *Move) TranslateGrsim() int {
	return 0
}

func (s *Stop) TranslateReal() []byte {
	message := make([]byte, 3)
	message[0] = 3
	message[1] = byte(RealStop)
	message[2] = byte(s.Id)

	return message
}

func (s *Stop) TranslateGrsim() int {
	return 0
}

func (k *Kick) TranslateReal() []byte {
	message := make([]byte, 4)
	message[0] = 4
	message[1] = byte(RealKick)
	message[2] = byte(k.Id)
	message[3] = byte(k.Speed)

	return message
}

func (i *Init) TranslateReal() []byte {
	message := make([]byte, 3)
	message[0] = 3
	message[1] = byte(RealInit)
	message[2] = byte(i.Id)

	return message
}

// ByteSplitInt splits a 32-bit integer into two bytes representing a 16-bit integer.
// The first return value is the upper byte, and the second is the lower byte.
// This function performs sign extension for negative numbers.
func ByteSplitInt16(n int) (byte, byte) {
	sixteenBitValue := int16(n)
	lowerByte := byte(sixteenBitValue & 0xFF)
	upperByte := byte((sixteenBitValue >> 8) & 0xFF)

	return upperByte, lowerByte
}

func float64To4Bytes(floatVal float64, bytes []byte) {
	floatBits := math.Float32bits(float32(floatVal))
	for i := 0; i < 4; i++ {
		bytes[i] = byte((floatBits >> (8 * (3 - i))) & 0xFF) // Extract bytes in big-endian order
	}
}

func appendPositionData(message []byte, startPos int, pos *mat.VecDense) int {
	// Process the first two elements as int16
	for i := 0; i < pos.Len()-1; i++ {
		upperByte, lowerByte := ByteSplitInt16(int(pos.AtVec(i)))
		message[startPos+i*2] = upperByte
		message[startPos+i*2+1] = lowerByte
	}

	float64To4Bytes(pos.AtVec(pos.Len()-1), message[startPos+(pos.Len()-1)*2:])

	return startPos + (pos.Len()-1)*2 + 4
}
