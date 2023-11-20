package gamestate

import "gonum.org/v1/gonum/mat"

// Enum type alias
//
// Defines type of shape of a line or arc
type FieldShape int32

const (
	Undefined FieldShape = iota
	CenterCircle
	TopTouchLine
	BottomTouchLine
	LeftGoalLine
	RightGoalLine
	HalfwayLine
	CenterLine
	LeftPenaltyStretch
	RightPenaltyStretch
	LeftFieldLeftPenaltyStretch
	LeftFieldRightPenaltyStretch
	RightFieldLeftPenaltyStretch
	RightFieldRightPenaltyStretch
)

type Field struct {
	// Field length (distance between goal lines) in mm
	FieldLengt int32

	// Field width (distance between touch lines) in mm
	FieldWidth int32

	// Field lines
	FieldLines []LineSegment

	// Field arcs
	FieldArcs []CircularArc

	// Goal width (distance between inner edges of goal posts) in mm
	GoalWidth int32

	// Goal depth (distance from outer goal line edge to inner goal back) in mm
	GoalDepth int32

	// Goal height in mm
	GoalHeight int32

	// Distance between the goal center and the center of the penalty mark in mm
	GoalToPenalty int32

	// Boundary width (distance from touch/goal line centers to boundary walls) in mm
	BoundaryWidth int32

	// Depth of the penalty/defense area (measured between line centers) in mm
	PenaltyAreaDepth int32

	// Width of the penalty/defense area (measured between line centers) in mm
	PenaltyAreaWidth int32

	// Radius of the center circle (measured between line centers) in mm
	CenterRadius int32

	// Thickness/width of the lines on the field in mm
	LineThickness int32

	// Ball radius in mm
	// (float type to represent sub-mm precision)
	BallRadius float32

	// Max allowed robot radius in mm
	// (float type to represent sub-mm precision)
	MaxRobotRadius float32
}

type LineSegment struct {
	// Name of marking
	Name string

	// Start point of line segment
	P1 *mat.VecDense

	// End point of line segment
	P2 *mat.VecDense

	// Thickness of line segment
	Thickness float32

	// Type of shape
	ShapeType FieldShape
}

type CircularArc struct {
	// Name of marking
	Name string

	// Center point of circular arc
	Center mat.VecDense

	// Radius of arc
	Radius float32

	// Start arngle in counter-clockwise order
	A1 float32

	// End angle in counter-clockwise order
	A2 float32

	// Thickness of arc
	Thickness float32

	// Type of shape
	ShapeType FieldShape
}

func (f *Field) addLine(
	name string,
	p1x float32,
	p1y float32,
	p2x float32,
	p2y float32,
	thickness float32,
	shape FieldShape) {

	p1 := []float64{float64(p1x), float64(p1y)}
	p2 := []float64{float64(p2x), float64(p2y)}
	line := LineSegment{
		Name:      name,
		P1:        mat.NewVecDense(2, p1),
		P2:        mat.NewVecDense(2, p2),
		Thickness: float32(thickness),
	}

	was := make([]int, 0)
}

// Check if Field contains some line
// with given name.
func (f *Field) hasLine(name string) bool {
	for _, line := range f.FieldLines {
		if line.Name == name {
			return true
		}
	}

	return false
}

// Check if Field contains some arc
// with given name.
func (f *Field) hasArc(name string) bool {
	for _, arc := range f.FieldArcs {
		if arc.Name == name {
			return true
		}
	}

	return false
}
