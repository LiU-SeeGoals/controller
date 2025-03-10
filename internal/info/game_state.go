package info

import (
	"encoding/json"
	"math"
	"time"

	. "github.com/LiU-SeeGoals/controller/internal/logger"
)

const TEAM_SIZE ID = 16

type GameState struct {
	Valid       bool
	Blue_team   *RobotTeam
	Yellow_team *RobotTeam

	Ball *Ball

	MessageReceived int64
}

// Constructor for game state
func NewGameState(capacity int) *GameState {
	gs := GameState{}
	gs.Valid = true

	gs.Ball = NewBall(capacity)
	var i ID
	gs.Blue_team = new(RobotTeam)
	gs.Yellow_team = new(RobotTeam)
	for i = 0; i < TEAM_SIZE; i++ {
		gs.Blue_team[i] = NewRobot(i, Blue, capacity)
		gs.Yellow_team[i] = NewRobot(i, Yellow, capacity)
	}
	return &gs
}

func (gs *GameState) ClosestRobot(target Position) (*Robot, float64) {
	shortestDistance := math.Inf(1)
	var closestRobot *Robot

	for _, robot := range gs.Blue_team {
		robotPos, err := robot.GetPosition()
		if err != nil {
			continue
		}

		distance := robotPos.Distance(target)
		if distance < shortestDistance {
			shortestDistance = distance
			closestRobot = robot
		}
	}
	for _, robot := range gs.Yellow_team {
		robotPos, err := robot.GetPosition()
		if err != nil {
			continue
		}

		distance := robotPos.Distance(target)
		if distance < shortestDistance {
			shortestDistance = distance
			closestRobot = robot
		}
	}

	return closestRobot, shortestDistance
}

// Should be called directly after inserting
// new data, recieved from ssl vision, into
// the game state
func (gs *GameState) Update() {

	// Update robot estimates
	// Has to be done before updating ball possessor
	for _, robot := range gs.Blue_team {
		if !robot.IsActive() {
			continue
		}
		// robot.Update()
	}
	for _, robot := range gs.Yellow_team {
		if !robot.IsActive() {
			continue
		}
		// robot.Update()
	}

	// Update ball estimate and possessor
	// Has to be done after updating robot estimates


	// --Ball possessor logic--
	// If new ball reading:
	//	estimate ball position
	//	estimate any possessor
	//	estimate if latest possessor has lost the ball
	// if no new ball reading:
	//	set ball position to position of latest possessor


	ballPos, recieved, err := gs.Ball.GetPositionTime()
	now := time.Now().UnixMilli()
	if err == nil { // There is a ball

		// If the timestamp of the latest ball position is
		// not old, we can check if any robot might be
		// in control of it, ie the possessor.
		if now-recieved < 100 { // WARN: Magic number, Position data is now-recieved ms old

			// Get the robot closest to the ball
			closestToBall, ballDistance := gs.ClosestRobot(ballPos)

			var facingBall bool
			if ballDistance > math.Inf(1) { // ballDistance will be inf if there is no robot on the field
				facingBall = false
			} else {
				facingBall = closestToBall.Facing(ballPos, 0.5) // WARN: Magic number
			}

			// I a robot is both facing the ball and in distance
			// consider it the possessor of the ball
			// TODO: Check matching velocities
			if ballDistance < 90 && facingBall { // WARN: Magic number, mm 
				gs.Ball.SetPossessor(closestToBall)

			// If the closest robot is to far away
			// set possessor to nil
			} else if (ballDistance > 110 ) { // WARN: Magic number, mm
				gs.Ball.SetPossessor(nil)

			}
		} else { 
			// Ball position to old to use, best we can do is assume the lateset
			// possessor still has it
			if gs.Ball.GetPossessor() != nil {
				// Set ball position to the position of the
				// robot which had the ball last
				possessorPos, _ := gs.Ball.GetPossessor().GetPosition()
				possessorPos.X += 90 * math.Cos(possessorPos.Angle) // WARN: Magic number
				possessorPos.Y += 90 * math.Sin(possessorPos.Angle) // WARN: Magic number
				gs.Ball.SetEstimatedPosition(possessorPos)
			}
		}
	} else {
		Logger.Errorf("Ball position retrieval failed: %v", err)
	}

}

// func (gs *GameState) ControlledBall() (bool, *Robot) {
//
// }
//
// func (gs *GameState) FreeBall() bool{
// 	
// }

func (gs *GameState) SetValid(valid bool) {
	gs.Valid = valid
}

func (gs *GameState) IsValid() bool {
	return gs.Valid
}

func (gs *GameState) SetYellowRobot(robotId uint32, x, y, angle float64, time int64) {
	gs.Yellow_team[robotId].SetPositionTime(x, y, angle, time)
}

func (gs *GameState) SetBlueRobot(robotId uint32, x, y, angle float64, time int64) {
	gs.Blue_team[robotId].SetPositionTime(x, y, angle, time)
}

// Updates position of robots and balls to their actual position
func (gs *GameState) SetBall(x, y, z float64, time int64) {
	gs.Ball.SetPositionTime(x, y, z, time)
}

func (gs *GameState) SetMessageReceivedTime(time int64) {
	gs.MessageReceived = time

}

func (gs *GameState) GetMessageReceivedTime() int64 {
	return gs.MessageReceived
}

func (gs *GameState) GetBall() *Ball {
	return gs.Ball
}

func (gs *GameState) GetTeam(team Team) *RobotTeam {
	if team == Yellow {
		return gs.Yellow_team
	} else {
		return gs.Blue_team
	}
}

func (gs *GameState) GetOtherTeam(team Team) *RobotTeam {
	if team != Yellow {
		return gs.Yellow_team
	} else {
		return gs.Blue_team
	}
}

func (gs *GameState) GetYellowRobots() *RobotTeam {
	return gs.Yellow_team
}

func (gs *GameState) GetBlueRobots() *RobotTeam {
	return gs.Blue_team
}

func (gs *GameState) GetRobot(id ID, team Team) *Robot {
	return gs.GetTeam(team)[id]
}

// String representation of game state
func (gs *GameState) String() string {
	gs_str := "{\n blue team: {\n"
	var i ID
	for i = 0; i < TEAM_SIZE; i++ {
		gs_str += "robot: {" + gs.Blue_team[i].String() + " },\n"
	}
	gs_str += "},\n yellow team: {\n"
	for i = 0; i < TEAM_SIZE; i++ {
		gs_str += "robot: {" + gs.Yellow_team[i].String() + " },\n"
	}
	// for _, line := range gs.Field.FieldLines {
	// 	gs_str += fmt.Sprintf("line: {%s}\n", line.String())
	// }
	// for _, arc := range gs.Field.FieldArcs {
	// 	gs_str += fmt.Sprintf("arc: {%s}\n", arc.String())
	// }
	gs_str += "}"
	return gs_str
}

type GameStateDTO struct {
	RobotPositions [2 * TEAM_SIZE]RobotDTO
	BallPosition   BallDTO
}

func (gs GameState) ToDTO() *GameStateDTO {
	gameStateDTO := GameStateDTO{}
	var i ID
	for i = 0; i < TEAM_SIZE; i++ {
		gameStateDTO.RobotPositions[i] = gs.GetRobot(i, Blue).ToDTO()
	}
	for i = 0; i < TEAM_SIZE; i++ {
		gameStateDTO.RobotPositions[TEAM_SIZE+i] = gs.GetRobot(i, Yellow).ToDTO()
	}
	gameStateDTO.BallPosition = gs.Ball.ToDTO()
	return &gameStateDTO
}

func (gs GameState) ToJson() []byte {
	gameStateJson, err := json.Marshal(gs.ToDTO())
	if err != nil {
		// fmt.Println("The gamestate packet could not be marshalled to JSON.")
		Logger.Error("The gamestate packet could not be marshalled to JSON.")
	}
	return gameStateJson
}
