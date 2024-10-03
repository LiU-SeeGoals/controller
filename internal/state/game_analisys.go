package state

type Zone struct {
	Scores []float32
}

type TeamAnalysis struct {
	Robots   [TEAM_SIZE]RobotAnalysis
	ZoneSize float32
	Zones    [][]Zone
}

type RobotAnalysis struct {
	Active           bool
	id               ID
	position         Position
	destination      Position
	velocity         Position
	MaxMoveSpeed     float32 // mm/s
	MaxRotationSpeed float32 // rad/s
	acceleration     float32 // mm/s^2
	deceleration     float32 // mm/s^2
}

type BallAnalysis struct {
	position    Position
	velocity    Position
	destination Position
}
type GameAnalysis struct {
	team      Team
	myTeam    TeamAnalysis
	otherTeam TeamAnalysis
	ball      BallAnalysis
}

type RobotAnalysisTeam [TEAM_SIZE]RobotAnalysis

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

func updateTeam(gameStateTeam *RobotTeam, teamAnalysis *TeamAnalysis) {
	for _, robot := range gameStateTeam {
		rAn := &teamAnalysis.Robots[robot.GetID()]
		if robot.IsActive() {
			rAn.Active = true
			rAn.id = robot.GetID()
			rAn.position = robot.GetPosition()
			rAn.velocity = robot.GetVelocity()

			if speed := calMoveSpeed(robot); speed > rAn.MaxMoveSpeed {
				rAn.MaxMoveSpeed = speed
			}
			if rotationSpeed := calRotationSpeed(robot); rotationSpeed > rAn.MaxRotationSpeed {
				rAn.MaxRotationSpeed = rotationSpeed
			}
			if acceleration := calAcceleration(robot); acceleration > rAn.acceleration {
				rAn.acceleration = acceleration
			}
			if deceleration := calDeceleration(robot); deceleration < rAn.deceleration {
				rAn.deceleration = deceleration
			}

		} else {
			rAn.Active = false
		}
	}
}

func updateBall(gameStateBall *Ball, ballAnalysis *BallAnalysis) {
	ballAnalysis.position = gameStateBall.GetPosition()
	ballAnalysis.velocity = gameStateBall.GetVelocity()
}

func NewGameAnalysis(fieldLength, fieldWidth, ZoneSize float32, team Team) *GameAnalysis {
	analysis := GameAnalysis{}
	hight := int(fieldLength / ZoneSize)
	width := int(fieldWidth / ZoneSize)
	analysis.team = team
	analysis.myTeam.ZoneSize = ZoneSize
	analysis.otherTeam.ZoneSize = ZoneSize
	analysis.myTeam.Zones = make([][]Zone, hight)
	analysis.otherTeam.Zones = make([][]Zone, hight)

	// Initialize the Zones
	for i := 0; i < hight; i++ {
		analysis.myTeam.Zones[i] = make([]Zone, width)
		analysis.otherTeam.Zones[i] = make([]Zone, width)
	}
	return &analysis
}

func (analysis *GameAnalysis) Update(gameState *GameState) {
	updateTeam(gameState.GetTeam(analysis.team), &analysis.myTeam)
	updateTeam(gameState.GetOtherTeam(analysis.team), &analysis.otherTeam)
	updateBall(gameState.GetBall(), &analysis.ball)
	// TODO: Update Zones
}
