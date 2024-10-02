package state

type Zone struct {
	scores []float32
}

type TeamAnalysis struct {
	Robots   [TEAM_SIZE]RobotAnalysis
	zoneSize float32
	zones    [][]Zone
}

type RobotAnalysis struct {
	active           bool
	id               ID
	position         Position
	destination      Position
	velocity         Position
	maxMoveSpeed     float32 // mm/s
	maxRotationSpeed float32 // rad/s
	acceleration     float32 // mm/s^2
	deceleration     float32 // mm/s^2
}
type GameAnalysis struct {
	team      Team
	myTeam    TeamAnalysis
	otherTeam TeamAnalysis
}

func calMoveSpeed(robot *Robot) float32 {
	velocity := robot.GetVelocity()
	return velocity.Norm()
}

func calRotationSpeed(robot *Robot) float32 {
	velocity := robot.GetVelocity()
	return velocity.Angel
}

func calAcceleration(robot *Robot) float32 {
	acceleration := robot.GetAcceleration()
	return acceleration
}

func calDeceleration(robot *Robot) float32 {
	return -calAcceleration(robot)
}

func updateTeam(gameStateTeam [TEAM_SIZE]*Robot, teamAnalysis *TeamAnalysis) {
	for _, robot := range gameStateTeam {
		rAn := &teamAnalysis.Robots[robot.GetID()]
		if robot.IsActive() {
			rAn.active = true
			rAn.id = robot.GetID()
			rAn.position = robot.GetPosition()
			rAn.velocity = robot.GetVelocity()

			if speed := calMoveSpeed(robot); speed > rAn.maxMoveSpeed {
				rAn.maxMoveSpeed = speed
			}
			if rotationSpeed := calRotationSpeed(robot); rotationSpeed > rAn.maxRotationSpeed {
				rAn.maxRotationSpeed = rotationSpeed
			}
			if acceleration := calAcceleration(robot); acceleration > rAn.acceleration {
				rAn.acceleration = acceleration
			}
			if deceleration := calDeceleration(robot); deceleration < rAn.deceleration {
				rAn.deceleration = deceleration
			}

		} else {
			rAn.active = false
		}
	}
}

func (analysis *GameAnalysis) Update(gameState *GameState, team Team) {
	zoneSize := 100
	updateTeam(gameState.GetTeam(team), &analysis.myTeam)
	updateTeam(gameState.GetOtherTeam(team), &analysis.otherTeam)
	analysis.team = team
	analysis.myTeam.zoneSize = float32(zoneSize)
	analysis.otherTeam.zoneSize = float32(zoneSize)

}
