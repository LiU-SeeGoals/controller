package state

type Zone struct {
	Scores []float32
	Score  float32
}

type TeamAnalysis struct {
	Robots   [TEAM_SIZE]RobotAnalysis
	ZoneSize float32
	Zones    [][]Zone
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

func (r *RobotAnalysis) IsActive() bool {
	return r.active
}

func (r *RobotAnalysis) GetID() ID {
	return r.id
}

func (r *RobotAnalysis) GetPosition() Position {
	return r.position
}

func (r *RobotAnalysis) GetDestination() Position {
	return r.destination
}

func (r *RobotAnalysis) GetVelocity() Position {
	return r.velocity
}

func (r *RobotAnalysis) GetMaxMoveSpeed() float32 {
	return r.maxMoveSpeed
}

func (r *RobotAnalysis) GetMaxRotationSpeed() float32 {
	return r.maxRotationSpeed
}

func (r *RobotAnalysis) GetAcceleration() float32 {
	return r.acceleration
}

func (r *RobotAnalysis) GetDeceleration() float32 {
	return r.deceleration
}

func (r *RobotAnalysis) SetDestination(destination Position) {
	r.destination = destination
}

type BallAnalysis struct {
	position    Position
	velocity    Position
	destination Position
}

func (b *BallAnalysis) GetPosition() Position {
	return b.position
}

func (b *BallAnalysis) GetVelocity() Position {
	return b.velocity
}

func (b *BallAnalysis) GetDestination() Position {
	return b.destination
}

func (b *BallAnalysis) SetDestination(destination Position) {
	b.destination = destination
}

type FieldInfo struct {
	Length float32
	Width  float32
}

type GameAnalysis struct {
	team      Team
	MyTeam    TeamAnalysis
	OtherTeam TeamAnalysis
	Ball      BallAnalysis
	FieldInfo FieldInfo
}

type RobotAnalysisTeam []RobotAnalysis

func calMoveSpeed(robot *Robot) float32 {
	velocity := robot.GetVelocity()
	return velocity.Norm()
}

func calRotationSpeed(robot *Robot) float32 {
	velocity := robot.GetVelocity()
	return velocity.Angle
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

func updateBall(gameStateBall *Ball, ballAnalysis *BallAnalysis) {
	ballAnalysis.position = gameStateBall.GetPosition()
	ballAnalysis.velocity = gameStateBall.GetVelocity()
}

func NewGameAnalysis(fieldLength, fieldWidth, zoneSize float32, team Team) *GameAnalysis {
	analysis := GameAnalysis{}
	hight := int(fieldLength / zoneSize)
	width := int(fieldWidth / zoneSize)
	analysis.team = team
	analysis.MyTeam.ZoneSize = zoneSize
	analysis.OtherTeam.ZoneSize = zoneSize
	analysis.FieldInfo.Length = fieldLength
	analysis.FieldInfo.Width = fieldWidth
	analysis.MyTeam.Zones = make([][]Zone, hight)
	analysis.OtherTeam.Zones = make([][]Zone, hight)

	// Initialize the Zones
	for i := 0; i < hight; i++ {
		analysis.MyTeam.Zones[i] = make([]Zone, width)
		analysis.OtherTeam.Zones[i] = make([]Zone, width)
	}
	return &analysis
}

func (analysis *GameAnalysis) UpdateState(gameState *GameState) {
	updateTeam(gameState.GetTeam(analysis.team), &analysis.MyTeam)
	updateTeam(gameState.GetOtherTeam(analysis.team), &analysis.OtherTeam)
	updateBall(gameState.GetBall(), &analysis.Ball)
}

func updateZone(team TeamAnalysis, fieldInfo FieldInfo, zoneSize float32, scoringFunc func(x float32, y float32, robots RobotAnalysisTeam) float32) {
	// Update the zones
	for i := 0; i < len(team.Zones); i++ {
		for j := 0; j < len(team.Zones[i]); j++ {
			// middle of the playing field in 0,0 so the zone need to be adjusted to the correct position
			x := float32(i)*zoneSize - fieldInfo.Length/2 + zoneSize/2
			y := float32(j)*zoneSize - fieldInfo.Width/2 + zoneSize/2
			team.Zones[i][j].Score = scoringFunc(x, y, team.Robots[:])
		}
	}
}

func (analysis *GameAnalysis) UpdateMyZones(scoringFunc func(x float32, y float32, robots RobotAnalysisTeam) float32) {
	updateZone(analysis.MyTeam, analysis.FieldInfo, analysis.MyTeam.ZoneSize, scoringFunc)
}

func (analysis *GameAnalysis) UpdateOtherZones(scoringFunc func(x float32, y float32, robots RobotAnalysisTeam) float32) {
	updateZone(analysis.OtherTeam, analysis.FieldInfo, analysis.OtherTeam.ZoneSize, scoringFunc)
}
