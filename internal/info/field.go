package info

// This package has been officially REdebloated
// ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⢤⣚⣟⠿⠯⠿⠷⣖⣤⠤⣀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⠀⢀⡠⡲⠟⠛⠉⠉⠀⠀⠀⠀⠀⠀⠀⠉⠓⠽⣢⣀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⣠⣔⠝⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⠤⣤⣤⡈⠹⣧⡄⠀⠀
// ⠀⠀⠀⢀⣴⠝⠁⠀⠀⠀⣴⣖⣚⣛⠽⠆⢀⠀⠀⠀⠙⠉⠉⠛⠁⠀⠈⢞⢆⠀
// ⠀⠀⢠⣻⠋⠀⠀⠀⠀⠀⠙⠋⠉⠀⠀⠀⠈⢣⠀⠈⡆⠀⠀⠀⠀⠀⠀⠀⢫⠆
// ⠀⢰⣳⣣⠔⠉⢱⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢇⠘⠀⠀⢀⣤⣶⣶⣶⣶⣾⣗
// ⢠⢯⠋⠀⠠⡴⠋⠙⠢⡄⠀⣠⡤⢔⣶⣶⣶⣶⣼⣤⣴⣾⡿⠋⢁⣤⣤⣽⣿⡏
// ⢸⣸⠀⠒⢻⣧⣶⣤⣴⣗⡏⠁⠀⠀⠀⠀⠀⠈⢻⣿⣿⣿⣠⣿⣿⣿⣿⣿⣿⠁
// ⣸⡏⠀⠘⠃⡿⢟⠇⢀⡿⣧⡄⠀⠀⠀⠀⠀⠀⣠⣿⠻⣿⣿⣿⣿⣿⣿⣿⠋⠀
// ⣷⠃⠀⠀⡇⡇⠀⣱⠞⠁⠸⣿⣦⡀⠀⠀⠀⠀⣸⠏⠀⠙⠻⢿⢿⣿⡟⠋⠀⢀
// ⢻⠀⠀⠀⣇⠴⠚⠁⠀⠀⠀⠈⠛⠿⢿⠤⠴⠚⠁⠀⣀⣠⠤⢔⡿⡟⠀⠀⠀⢸
// ⣇⣘⡓⣺⡿⠀⠀⢠⠶⠒⢶⣲⡒⠒⠒⠒⠒⣛⣉⡩⠤⠖⠚⢁⡝⢠⡄⠀⢀⠦
// ⠙⢶⢏⠁⠀⠀⠀⠀⠀⠀⠀⠈⠙⠿⣟⡛⠉⠀⠀⠀⠀⢀⡤⠊⢀⡜⢀⡼⡸⠏
// ⠀⠀⢯⣦⠀⠀⠀⠀⠀⠀⠀⠀⢀⡀⠀⠉⠉⠓⠒⠚⠉⠁⠀⣠⠎⢠⡾⡽⠁⠀
// ⠀⠀⠈⠪⣵⠀⠀⠀⠀⠀⠀⠀⠀⠉⠳⠶⣤⣤⣤⣤⣤⡶⠟⣅⣴⣏⠏⠀⠀⠀
// ⠀⠀⠀⠀⠉⢳⣄⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣤⣴⡯⠋⠁⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⠀⠉⠢⣤⣄⣀⠀⠀⠀⠀⠀⠀⢀⣀⠮⠓⠉⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠈⠛⠓⠂⠀⠂⠁⠉⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀



import (
	"fmt"
)

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

// Holds all field geometry data
//
// Definitions are similar to the ones in the SSL ptoto files,
// but slight changes have been made for readability.
// Please compare with the generated proto files if
// anything is unclear.
type GameField struct {
	// Field length (distance between goal lines) in mm
	FieldLength int32

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
	// Is never set!
	// Goal height in mm
	GoalHeight int32
	// Is never set!
	// Distance between the goal center and the center of the penalty mark in mm
	GoalToPenalty int32

	// Boundary width (distance from touch/goal line centers to boundary walls) in mm
	BoundaryWidth int32

	// Depth of the penalty/defense area (measured between line centers) in mm
	PenaltyAreaDepth int32

	// Width of the penalty/defense area (measured between line centers) in mm
	PenaltyAreaWidth int32

	// Is never set!
	// Radius of the center circle (measured between line centers) in mm
	CenterRadius int32

	// Is never set!
	// Thickness/width of the lines on the field in mm
	LineThickness int32

	// Is never set!
	// Ball radius in mm
	// (float type to represent sub-mm precision)
	BallRadius float64

	// Is never set!
	// Max allowed robot radius in mm
	// (float type to represent sub-mm precision)
	MaxRobotRadius float64
}

func NewGameField() *GameField {
	return &GameField{}
}

func (gf *GameField) SetField(length, width, goalWidth, goalDepth, boundaryWidth, penaltyWidth, penaltyDepth int32) {
	gf.FieldLength = length
	gf.FieldWidth = width
	gf.GoalWidth = goalWidth
	gf.GoalDepth = goalDepth
	gf.BoundaryWidth = boundaryWidth
	gf.PenaltyAreaWidth = penaltyWidth
	gf.PenaltyAreaDepth = penaltyDepth
}
func (gf *GameField) AddFieldLine(name string, x1, y1, x2, y2, thickness float64, lineType int) {
	gf.SetLine(name, x1, y1, x2, y2, thickness, FieldShape(lineType))
}

func (gf *GameField) AddFieldArc(name string, centerX, centerY, radius, angle1, angle2, thickness float64, shape int) {
	gf.SetArc(name, centerX, centerY, radius, angle1, angle2, thickness, FieldShape(shape))
}

type Point struct {
	X float64
	Y float64
}

// Holds line segment data
//
// Shape enum maps one to one
// with the SSL vision enum
type LineSegment struct {
	// Name of marking
	Name string

	// Start point of line segment
	P1 Point

	// End point of line segment
	P2 Point

	// Thickness of line segment
	Thickness float64

	// Type of shape
	ShapeType FieldShape
}

// Holds arc data
//
// Shape enum maps one to one
// with the SSL vision enum
type CircularArc struct {
	// Name of marking
	Name string

	// Center point of circular arc
	Center Point

	// Radius of arc
	Radius float64

	// Start arngle in counter-clockwise order
	A1 float64

	// End angle in counter-clockwise order
	A2 float64

	// Thickness of arc
	Thickness float64

	// Type of shape
	ShapeType FieldShape
}

// Add a new line segment to Field object
func (f *GameField) SetLine(
	name string,
	p1x float64,
	p1y float64,
	p2x float64,
	p2y float64,
	thickness float64,
	shape FieldShape) {

	line := LineSegment{
		Name:      name,
		P1:        Point{X: p1x, Y: p1y},
		P2:        Point{X: p2x, Y: p2y},
		Thickness: float64(thickness),
	}

	f.FieldLines = append(f.FieldLines, line)
}

// Adds a new arc to Field object
func (f *GameField) SetArc(
	name string,
	centerX float64,
	centerY float64,
	radius float64,
	angle1 float64,
	angle2 float64,
	thickness float64,
	shape FieldShape) {

	arc := CircularArc{
		Name:      name,
		Center:    Point{X: centerX, Y: centerY},
		Radius:    radius,
		A1:        angle1,
		A2:        angle2,
		Thickness: thickness,
		ShapeType: shape,
	}
	f.FieldArcs = append(f.FieldArcs, arc)
}

// String representation of LineSegment
func (l *LineSegment) String() string {
	x1 := l.P1.X
	y1 := l.P1.Y

	x2 := l.P2.X
	y2 := l.P2.Y
	return fmt.Sprintf("name: %s, p1: {%f, %f}, p2: {%f, %f}", l.Name, x1, y1, x2, y2)
}

// String representation of CircularArc
func (a *CircularArc) String() string {
	x := a.Center.X
	y := a.Center.Y

	return fmt.Sprintf("name: %s, center: {%f, %f}, rad: %f, a1: %f, a2: %f", a.Name, x, y, a.Radius, a.A1, a.A2)
}
