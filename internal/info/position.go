package info

import (
	"fmt"
	"math"
)

type Position struct {
	X     float64
	Y     float64
	Z     float64
	Angle float64
}

func (p *Position) FacingPosition(target Position, threshold float64) bool {

	dx := target.X - p.X
	dy := target.Y - p.Y

	targetDirection := math.Atan2(float64(dy), float64(dx))
	currentDirection := p.Angle

	angleDiff := math.Abs(float64(targetDirection - currentDirection))
	if angleDiff < threshold {
		return true
	} else {
		return false
	}
}

func (p Position) AngleToPosition(p2 Position) float64 {
	dx := p2.X - p.X
	dy := p2.Y - p.Y
	return float64(math.Atan2(float64(dy), float64(dx)))
}

func (p Position) AngleDistance(p2 Position) float64 {
	diff := p.AngleToPosition(p2) - p.Angle
	return math.Abs(math.Remainder(diff, 2*math.Pi))
}

//Disntance between two points
func (p Position) Distance(p2 Position) float64 {
	dx := p.X - p2.X
	dy := p.Y - p2.Y
	dz := p.Z - p2.Z
	return float64(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
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

func (p Position) Dot(other Position) float64 {
	return p.X*other.X + p.Y*other.Y + p.Z*other.Z
}

func (p Position) Norm() float64 {
	return float64(math.Sqrt(float64(p.X*p.X + p.Y*p.Y + p.Z*p.Z)))
}

func (p Position) Scale(scalar float64) Position {
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
