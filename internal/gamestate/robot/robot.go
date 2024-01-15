// Package robot defines the structure and behavior of a robot within a game state.
package robot

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

// Robot represents a robot in the game state.
type Robot struct {
	id   int         // unique identifier for the robot
	team Team        // team to which the robot belongs
	pos  *mat.VecDense // position {x (mm), y (mm), angle (rad)}
	vel  *mat.VecDense // velocity {x (mm/s), y (mm/s), angular (rad/s)}
}

// NewRobot creates a new Robot with the given id and team.
// The initial position and velocity of the Robot are set to (0, 0, 0) in millimeters.
func NewRobot(id int, team Team) *Robot {
	return &Robot{
		id:   id,
		team: team,
		pos:  mat.NewVecDense(3, []float64{0, 0, 0}), // in mm
		vel:  mat.NewVecDense(3, []float64{0, 0, 0}), // in mm/s
	}
}

// GetId returns the ID of the robot.
func (r *Robot) GetId() int {
	return r.id
}

// GetTeam returns the team of the robot.
func (r *Robot) GetTeam() Team {
	return r.team
}

// GetPosition returns a pointer to a mat.VecDense representing the current position of the robot.
// The position vector contains three elements: X, Y coordinates in millimeters, and the orientation angle W in radians.
func (r *Robot) GetPosition() *mat.VecDense {
	return r.pos
}

// GetVelocity returns a pointer to a mat.VecDense representing the current velocity of the robot.
// The velocity vector contains three elements: velocity along the X and Y axis in millimeters per second,
// and the angular velocity W in radians per second.
func (r *Robot) GetVelocity() *mat.VecDense {
	return r.vel
}

// GetPositionX returns the X-coordinate of the robot's position in millimeters (mm).
func (r *Robot) GetPositionX() float64 {
	return r.pos.AtVec(0)
}

// GetPositionY returns the Y-coordinate of the robot's position in millimeters (mm).
func (r *Robot) GetPositionY() float64 {
	return r.pos.AtVec(1)
}

// GetAngle returns the angle of the robot in radians.
func (r *Robot) GetAngle() float64 {
	return r.pos.AtVec(2)
}

// GetVelocityX returns the velocity of the robot along the X-axis in millimeters per second (mm/s).
func (r *Robot) GetVelocityX() float64 {
	return r.vel.AtVec(0)
}

// GetVelocityY returns the velocity of the robot along the Y-axis in millimeters per second (mm/s).
func (r *Robot) GetVelocityY() float64 {
	return r.vel.AtVec(1)
}

// GetAngularVelocity returns the angular velocity of the robot in radians per second (rad/s).
func (r *Robot) GetAngularVelocity() float64 {
	return r.vel.AtVec(2)
}

// SetPosition updates the robot's current position to the specified X, Y coordinates in millimeters,
// and orientation angle in radians.
func (r *Robot) SetPosition(x, y, angle float64) {
	r.pos.SetVec(0, x)
	r.pos.SetVec(1, y)
	r.pos.SetVec(2, angle)
}

// SetVelocity updates the robot's current velocity to the specified components along the X and Y axes in
// millimeters per second, and angular velocity in radians per second.
func (r *Robot) SetVelocity(x_velocity, y_velocity, angular_velocity float64) {
	r.vel.SetVec(0, x_velocity)
	r.vel.SetVec(1, y_velocity)
	r.vel.SetVec(2, angular_velocity)
}

// String provides a formatted string representation of the robot's current state, including its ID,
// position, and velocity. The position and velocity are given as tuples of X, Y, and orientation angle
// or angular velocity, respectively.
func (r *Robot) String() string {
	x := r.GetPositionX();
	y := r.GetPositionY();
	angle := r.GetAngle();

	vX := r.GetVelocityX();
	vY := r.GetVelocityY();
	angularVelocity := r.GetAngularVelocity();

	posString := fmt.Sprintf("(%f, %f, %f)", x, y, angle)
	velString := fmt.Sprintf("(%f, %f, %f)", vX, vY, angularVelocity)

	return fmt.Sprintf("Robot{id: %d, pos: %s, vel: %s}", r.id, posString, velString)
}
