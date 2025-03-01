package ai

import (
	// "fmt"
	"fmt"
	"math"

	"github.com/LiU-SeeGoals/controller/internal/info"
)

type MovementComposition struct{}

const ( // AvoidCollision constants
	detectionRadius = float32(1000.0) // How close is "too close" to another robot (mm)
	baseDamping     = float32(50.0)     // baseline damping
	brakeDampingMax = float32(1000.0)    // max damping when very close
	wGyroMax        = float32(0.0) // max gyroscopic gain
	stepClamp       = float32(2000.0) // maximum allowed step (mm)
)

const ( // SimpleAvoid constants
	kAttraction  = float32(1.0)   // how strongly we pull toward the goal
	kRepulsion   = float32(10.0) // how strongly we push away from nearby obstacles
	detectionRad = float32(500.0) // detection radius for obstacles (mm)
)

const ( // Juggernaut constants
	detectionRadJuggernaut = float32(500.0) // detection radius for obstacles (mm)
)

// Dont care about obstacles being in the way, 
// just slows down do avoid penalties.
func (mc *MovementComposition) Juggernaut(robot *info.Robot, goal info.Position, gs *info.GameState) info.Position {
	// Current state
	obstacles := getObstacles(gs, robot)

	slowDown := false
	for _, obs := range obstacles {
		dist := robot.GetPosition().Distance(obs)
		if dist < detectionRadJuggernaut {
			slowDown = true
			break
		}
	}

	if slowDown {
		pos := goal.Sub(robot.GetPosition()).Normalize().Scale(300)
		return pos
	}

	return goal
}


// SimpleAvoid: Just a force pulling the robot toward a goal,
// and a repulsive force pushing away from obstacles within detectionRad.
// Good for avoiding contact with other robots, can get stuck in local minima though.
func (mc *MovementComposition) SimpleAvoid(robot *info.Robot, goal info.Position, gs *info.GameState) info.Position {
	// Current state
	obstacles := getObstacles(gs, robot)
	pos := robot.GetPosition()

	// 1) Attractive force
	Fatt := computeAttractiveForce(pos, goal, kAttraction)

	// 2) Repulsive force from obstacles
	Frep := computeRepulsiveForce(pos, obstacles, kRepulsion, detectionRad)

	// 3) Sum the forces => acceleration (mass = 1)
	Ftotal := Fatt.Add(Frep)

	newPos := pos.Add(Ftotal)

	return newPos
}

// computeAttractiveForce: -k*(pos - goal), i.e. a vector from pos→goal
func computeAttractiveForce(pos, goal info.Position, k float32) info.Position {
	dx := goal.X - pos.X
	dy := goal.Y - pos.Y
	return info.Position{
		X: k * dx,
		Y: k * dy,
	}
}

func computeRepulsiveForce(pos info.Position, obstacles []info.Position, k, rDet float32) info.Position {
	force := info.Position{X: 0, Y: 0}
	for _, obs := range obstacles {
		dx := pos.X - obs.X
		dy := pos.Y - obs.Y

		dist := float32(math.Sqrt(float64(dx*dx + dy*dy)))
		fmt.Println("Distance: ", dist)

		if dist < rDet && dist != 0{

			// push away
			mag := k * (rDet - dist)
			// direction: from obs → pos
			force.X += mag * (dx / dist)
			force.Y += mag * (dy / dist)
		}
	}
	return force
}

func getObstacles(gs *info.GameState, robot *info.Robot) []info.Position {

	var obstacles []info.Position

	// Handle own team avoiding
	for _, otherRobot := range gs.Yellow_team {
		if !otherRobot.IsActive() { continue }
		// Avoid self
		if otherRobot.GetID() != robot.GetID() {
			obstacles = append(obstacles, otherRobot.GetPosition())

		}
	}

	for _, otherRobot := range gs.Blue_team {
		if !otherRobot.IsActive() { continue }
		obstacles = append(obstacles, otherRobot.GetPosition())
	}

	return obstacles
}

// AvoidCollision computes a single-step “waypoint” that points the robot
// in the direction of the sum of potential, gyroscopic, and damping/braking forces.
// Good for moving in a swarm with friendly robots.
func (mc *MovementComposition) AvoidCollision(robot *info.Robot, goal info.Position, gs *info.GameState) info.Position {
	// Current position and velocity
	robotPos := robot.GetPosition()
	oldVel := robot.GetVelocity()

	// 1) Find nearest obstacle/robot
	nearestRobot := getNearestRobot(gs, robot)

	// 2) Potential force (robot → goal)
	Fp := computePotentialForce(robotPos, goal)

	// 3) Gyroscopic force
	Fg := computeGyroscopicForce(robotPos, oldVel, nearestRobot)

	// 4) Damping + Braking force
	Fd := computeDampingForce(oldVel, robotPos, nearestRobot)

	// Sum up the three terms: total “direction” we want to move
	totalForce := Fp.Add(Fg)
	totalForce = totalForce.Add(Fd)

	// (Optional) Clamp the step size so we don’t send a giant jump
	stepSize := totalForce.Norm()
	if stepSize > stepClamp {
		totalForce = totalForce.Scale(stepClamp / stepSize)
	}

	// This code returns a single step from current position
	// in the direction of totalForce
	newPos := robotPos.Add(totalForce)

	// Return newPos as the waypoint
	return newPos
}

// computePotentialForce: points from robot to goal, optionally magnitude-clamped
func computePotentialForce(robotPos, goal info.Position) info.Position {
	attractive := goal.Sub(robotPos)
	mag := attractive.Norm()
	maxMag := float32(2000.0)
	if mag > maxMag {
		attractive = attractive.Scale(maxMag / mag)
	}
	return attractive
}

// computeGyroscopicForce: rotates velocity ±90° depending on (obsVec × vel),
// scaled by how close the other robot is and how large the velocity is.
func computeGyroscopicForce(robotPos, vel info.Position, other *info.Robot) info.Position {
	if other == nil {
		return info.Position{X: 0, Y: 0}
	}

	// Vector from us to nearest robot
	otherPos := other.GetPosition()
	obsVec := otherPos.Sub(robotPos)

	// Speed
	speed := vel.Norm()
	if speed < 1e-8 {
		return info.Position{X: 0, Y: 0}
	}

	// Cross sign
	crossVal := obsVec.Cross2D(vel) // 2D cross => obsVec.X*vel.Y - obsVec.Y*vel.X
	signDir := float32(1.0)
	if crossVal < 0 {
		signDir = -1.0
	}

	// ±90 deg rotation of velocity
	vPerp := info.Position{
		X: -vel.Y * signDir,
		Y: vel.X * signDir,
	}

	// Distance-based scale => stronger if close
	dist := obsVec.Norm()
	if dist < 1e-3 {
		dist = 1e-3
	}
	scale := float32(0.0)
	if dist < detectionRadius {
		scale = (detectionRadius - dist) / detectionRadius
	}
	magnitude := wGyroMax * scale * speed

	// Normalize vPerp => multiply by magnitude
	perpLen := vPerp.Norm()
	if perpLen < 1e-8 {
		return info.Position{X: 0, Y: 0}
	}
	return vPerp.Scale(magnitude / perpLen)
}

// computeDampingForce: base damping plus extra “brake” if the other robot is close
func computeDampingForce(vel, robotPos info.Position, other *info.Robot) info.Position {
	if other == nil {
		return vel.Scale(-baseDamping)
	}

	otherPos := other.GetPosition()
	distance := robotPos.Distance(otherPos)

	if distance >= detectionRadius {
		// beyond detection => base damping only
		return vel.Scale(-baseDamping)
	}

	ratio := (detectionRadius - distance) / detectionRadius
	D := baseDamping + ratio*(brakeDampingMax-baseDamping)
	return vel.Scale(-D)
}

// getNearestRobot scans both teams for the closest robot under detectionRadius, ignoring self.
func getNearestRobot(gs *info.GameState, robot *info.Robot) *info.Robot {
	minDist := float32(math.MaxFloat32)
	var nearest *info.Robot

	robotPos := robot.GetPosition()

	// check same team
	for _, other := range gs.GetTeam(robot.GetTeam()) {
		if other.GetID() == robot.GetID() { continue }
		if !other.IsActive() { continue }

		otherPos := other.GetPosition()
		dist := robotPos.Distance(otherPos)
		if dist < minDist && dist < detectionRadius {
			minDist = dist
			nearest = other
		}
	}
	// check other team
	for _, other := range gs.GetOtherTeam(robot.GetTeam()) {
		if !other.IsActive() { continue }

		otherPos := other.GetPosition()
		dist := robotPos.Distance(otherPos)
		if dist < minDist && dist < detectionRadius {
			minDist = dist
			nearest = other
		}
	}
	return nearest
}
