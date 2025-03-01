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

func (p Position) String() string {
	return fmt.Sprintf("(%f, %f, %f, %f)", p.X, p.Y, p.Z, p.Angle)
}

// Cross2D calculates the cross product of two 2D vectors
func (p Position) Cross2D(other Position) float32 {
    return p.X*other.Y - p.Y*other.X
}

// Distance calculates the Euclidean distance between two positions
func (p Position) Distance(other Position) float32 {
    dx := float64(p.X - other.X)
    dy := float64(p.Y - other.Y)
    return float32(math.Sqrt(dx*dx + dy*dy))
}

func (p Position) Add(other Position) Position {
	return Position{
		X:     p.X + other.X,
		Y:     p.Y + other.Y,
		Z:     p.Z + other.Z,
		Angle: p.Angle + other.Angle,
	}
}

func (p Position) Sub(other Position) Position {
	return Position{
		X:     p.X - other.X,
		Y:     p.Y - other.Y,
		Z:     p.Z - other.Z,
		Angle: p.Angle - other.Angle,
	}
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
