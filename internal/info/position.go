package info

import (
	"fmt"
	"math"
)

type Position struct {
	X     float32
	Y     float32
	Z     float32
	Angle float32
}

func (p Position) AngleToPosition(p2 Position) float32 {
	dx := p2.X - p.X
	dy := p2.Y - p.Y
	return float32(math.Atan2(float64(dy), float64(dx)))
}

//Disntance between two points
func (p Position) Distance(p2 Position) float32 {
	dx := p.X - p2.X
	dy := p.Y - p2.Y
	dz := p.Z - p2.Z
	return float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
}

func (p Position) String() string {
	return fmt.Sprintf("(%f, %f, %f, %f)", p.X, p.Y, p.Z, p.Angle)
}

func (p Position) Add(other *Position) Position {
	return Position{
		X:     p.X + other.X,
		Y:     p.Y + other.Y,
		Z:     p.Z + other.Z,
		Angle: p.Angle + other.Angle,
	}
}

func (p Position) Sub(other *Position) Position {
	return Position{
		X:     p.X - other.X,
		Y:     p.Y - other.Y,
		Z:     p.Z - other.Z,
		Angle: p.Angle - other.Angle,
	}
}

func (p Position) Dot(other Position) float32 {
	return p.X*other.X + p.Y*other.Y + p.Z*other.Z
}

func (p Position) Norm() float32 {
	return float32(math.Sqrt(float64(p.X*p.X + p.Y*p.Y + p.Z*p.Z)))
}

func (p Position) Scale(scalar float32) Position {
	return Position{
		X:     p.X * scalar,
		Y:     p.Y * scalar,
		Z:     p.Z * scalar,
		Angle: p.Angle * scalar,
	}
}

func (p Position) Normalize() Position {
	norm := p.Norm()
	return Position{
		X:     p.X / norm,
		Y:     p.Y / norm,
		Z:     p.Z / norm,
		Angle: p.Angle / norm,
	}
}

func (p Position) ToDTO() string {
	return fmt.Sprintf("(%f, %f, %f, %f)", p.X, p.Y, p.Z, p.Angle)
}
