package ai

// import (
// 	"fmt"
// 	"math"
// 	"time"
//
// 	"github.com/LiU-SeeGoals/controller/internal/info"
// )
//
// const (
// 	RUNNING int = iota
// 	COMPLETE
// 	TIME_EXPIRED
// 	ERROR
// 	FAILED
// )
//
// type ScenarioSlowBrain struct {
// 	team              info.Team
// 	incomingGameState <-chan info.GameState
// 	outgoingPlan      chan<- info.GamePlan
// 	scale             float32
// 	run_scenario      int // -1 for all
// }
//
// func NewScenarioSlowBrain(scale float32, run_scenario int) *ScenarioSlowBrain {
// 	return &ScenarioSlowBrain{scale: scale, run_scenario: run_scenario}
// }
//
// func (sb *ScenarioSlowBrain) Init(incoming <-chan info.GameState, outgoing chan<- info.GamePlan, team info.Team) {
// 	sb.incomingGameState = incoming
// 	sb.outgoingPlan = outgoing
// 	sb.team = team
//
// 	go sb.Run()
// }
//
// type ScenarioTest interface {
// 	Run() []*info.Instruction
// 	Archived(*info.GameState) int
// }
//
// func (sb ScenarioSlowBrain) Run() {
// 	var gameState info.GameState
// 	gameState.SetValid(false)
//
// 	scenarios := []ScenarioTest{}
// 	scenarios = append(scenarios, NewMoveToTest(sb.team))
// 	scenarios = append(scenarios, NewMoveToBallTest(sb.team))
// 	scenarios = append(scenarios, NewObstacleAvoidanceTest(sb.team))
// 	scenarios = append(scenarios, NewRealTest(sb.team))
// 	scenarios = append(scenarios, NewMoveWithBallToPositionTest(sb.team))
// 	scenarios = append(scenarios, NewKickToPlayerTest(sb.team))
// 	// scenarios = append(scenarios, NewObstacleAvoidanceTest(sb.team))
//
// 	scenario_index := 0
// 	if sb.run_scenario >= 0 {
// 		scenario_index = sb.run_scenario
// 	}
//
// 	fmt.Println("Running scenarios")
// 	for {
// 		gameState = <-sb.incomingGameState
//
// 		if !gameState.IsValid() {
// 			fmt.Println("ScenarioSlowBrain: Invalid game state")
// 			time.Sleep(40 * time.Millisecond)
// 			continue
// 		}
//
// 		scenario := scenarios[scenario_index]
// 		game_state := scenario.Archived(&gameState)
// 		if game_state == COMPLETE {
// 			fmt.Println("Scenario", scenario_index, "completed")
// 		} else if game_state == TIME_EXPIRED {
// 			fmt.Println("Scenario", scenario_index, "time expired")
// 		} else if game_state == FAILED {
// 			fmt.Println("Scenario", scenario_index, "failed")
// 		}
// 		if game_state != RUNNING {
// 			scenario_index++
// 			if scenario_index >= len(scenarios) || sb.run_scenario >= 0 {
// 				panic("ScenarioSlowBrain: No more scenarios") // TODO: Handle this better
// 			}
// 			scenario = scenarios[scenario_index]
// 		}
//
// 		inst := scenario.Run()
//
// 		plan := info.GamePlan{}
// 		plan.Instructions = inst
//
// 		plan.Team = sb.team
//
// 		plan.Valid = true
//
// 		sb.outgoingPlan <- plan
//
// 	}
// }
//
// type MoveToBallTest struct {
// 	team     info.Team
// 	at_state int
// 	start    time.Time
// 	max_time time.Duration
// }
//
// func NewMoveToBallTest(team info.Team) *MoveToBallTest {
// 	return &MoveToBallTest{
// 		team:     team,
// 		max_time: 30 * time.Second,
// 		at_state: -1,
// 	}
// }
//
// func (m *MoveToBallTest) Run() []*info.Instruction {
// 	if m.at_state == -1 {
// 		m.start = time.Now()
// 		m.at_state = 0
// 	}
// 	if m.at_state == 0 {
// 		return []*info.Instruction{
// 			{Type: info.MoveToBall, Id: 0},
// 		}
// 	} else {
// 		return []*info.Instruction{
// 			{Type: info.MoveToBall, Id: 0},
// 		}
// 	}
// }
//
// func (m *MoveToBallTest) Archived(gs *info.GameState) int {
// 	robot_pos := gs.GetRobot(info.ID(0), m.team).GetPosition()
// 	ball_pos := gs.GetBall().GetPosition()
//
// 	dx := float64(robot_pos.X - ball_pos.X)
// 	dy := float64(robot_pos.Y - ball_pos.Y)
// 	distance := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
//
// 	if m.at_state == 0 {
// 		if distance < 500 {
// 			m.at_state = 1
// 		}
// 	} else if m.at_state == 1 {
// 		if distance > 500 {
// 			fmt.Println("Failed with robot at (", robot_pos.X, robot_pos.Y, ") and ball at (", ball_pos.X, ball_pos.Y, ")")
// 			m.at_state = 2
// 		}
// 	}
//
// 	if m.at_state >= 0 {
// 		if time.Since(m.start) > m.max_time || m.at_state == 2 {
// 			if m.at_state == 0 {
// 				fmt.Println("Did not reach ball")
// 				return TIME_EXPIRED
// 			} else if m.at_state == 1 {
// 				fmt.Println("Reached ball and stayed there! :D")
// 				return COMPLETE
// 			} else {
// 				fmt.Println("Reached ball but then lost it :(")
// 				return FAILED
// 			}
// 		}
// 	}
// 	return RUNNING
// }
//
// // --------------------------Real Test--------------------------------
// type RealTest struct {
// 	team           info.Team
// 	at_state       int
// 	start          time.Time
// 	max_time       time.Duration
// 	instructionSet []*info.Instruction
// }
//
// func NewRealTest(team info.Team) *RealTest {
//
// 	instSet := []*info.Instruction{
// 		{Type: info.MoveToPosition, Id: 5, Position: info.Position{X: -2000, Y: 0}},
// 		{Type: info.MoveToPosition, Id: 5, Position: info.Position{X: 2000, Y: 0}},
// 	}
//
// 	return &RealTest{
// 		team:           team,
// 		max_time:       600 * time.Second,
// 		at_state:       0,
// 		instructionSet: instSet,
// 	}
// }
//
// func (m *RealTest) Run() []*info.Instruction {
//
// 	return []*info.Instruction{m.instructionSet[m.at_state]}
// }
//
// func (m *RealTest) Archived(gs *info.GameState) int {
// 	// if len(gs.Yellow_team) != 4 {
// 	// 	fmt.Println("Yellow team has", len(gs.Yellow_team), "4 required")
// 	// 	return ERROR
// 	// } else if len(gs.Blue_team) != 3 {
// 	// 	fmt.Println("Blue team has", len(gs.Blue_team), "3 required")
// 	// 	return ERROR
// 	// }
//
// 	// This assumes that the robot ids range from 0 to m.team_size
// 	target := m.instructionSet[m.at_state].Position
// 	robot := gs.GetRobot(info.ID(5), m.team)
//
// 	if atPosition(robot, target) {
// 		m.at_state = (m.at_state + 1) % 2 // Cycle through the states
// 	}
//
// 	return RUNNING
// }
//
// // --------------------------Obstacle avoidance--------------------------------
// type ObstacleAvoidanceTest struct {
// 	team           info.Team
// 	team_size      int
// 	at_states      []int
// 	start          time.Time
// 	max_time       time.Duration
// 	instructionSet [][]info.Instruction
// }
//
// func NewObstacleAvoidanceTest(team info.Team) *ObstacleAvoidanceTest {
// 	team_size := 4
// 	var instSet [][]info.Instruction
// 	if team == info.Yellow {
// 		instSet = [][]info.Instruction{
// 			{
// 				{Type: info.MoveToPosition, Id: 0, Position: info.Position{X: -4000, Y: 0}},
// 				{Type: info.MoveToPosition, Id: 0, Position: info.Position{X: 4000, Y: 0}},
// 			},
// 			{
// 				{Type: info.MoveToPosition, Id: 1, Position: info.Position{X: 640, Y: -1000}},
// 				{Type: info.MoveToPosition, Id: 1, Position: info.Position{X: 640, Y: 1000}},
// 			},
// 			{ // --------------------------Obstacle avoidance--------------------------------
// 				{Type: info.MoveToPosition, Id: 2, Position: info.Position{X: 1925, Y: -1000}},
// 				{Type: info.MoveToPosition, Id: 2, Position: info.Position{X: 1925, Y: 1000}},
// 			},
// 			{
//
// 				{Type: info.MoveToPosition, Id: 3, Position: info.Position{X: 3210, Y: -1000}},
// 				{Type: info.MoveToPosition, Id: 3, Position: info.Position{X: 3210, Y: 1000}},
// 			},
// 			{
//
// 				{Type: info.MoveToPosition, Id: 3, Position: info.Position{X: 3210, Y: -1000}},
// 				{Type: info.MoveToPosition, Id: 3, Position: info.Position{X: 3210, Y: 1000}},
// 			},
// 		}
// 	} else {
// 		instSet = [][]info.Instruction{
// 			{
// 				{Type: info.MoveToPosition, Id: 0, Position: info.Position{X: -3215, Y: 1000}},
// 				{Type: info.MoveToPosition, Id: 0, Position: info.Position{X: -3215, Y: -1000}},
// 			},
// 			{
// 				{Type: info.MoveToPosition, Id: 1, Position: info.Position{X: -1930, Y: -1000}},
// 				{Type: info.MoveToPosition, Id: 1, Position: info.Position{X: -1930, Y: 1000}},
// 			},
// 			{
// 				{Type: info.MoveToPosition, Id: 2, Position: info.Position{X: -645, Y: -1000}},
// 				{Type: info.MoveToPosition, Id: 2, Position: info.Position{X: -645, Y: 1000}},
// 			}}
// 	}
// 	nr_robots := len(instSet)
// 	return &ObstacleAvoidanceTest{
// 		team:           team,
// 		team_size:      team_size,
// 		max_time:       600 * time.Second,
// 		at_states:      make([]int, nr_robots),
// 		instructionSet: instSet,
// 	}
// }
//
// func (m *ObstacleAvoidanceTest) Run() []*info.Instruction {
// 	var instructions []*info.Instruction
//
// 	// This assumes that the robot ids range from 0 to m.team_size
// 	for id, at_state := range m.at_states {
// 		instructions = append(instructions, &m.instructionSet[id][at_state])
// 	}
// 	return instructions
// }
//
// func (m *ObstacleAvoidanceTest) Archived(gs *info.GameState) int {
// 	// if len(gs.Yellow_team) != 4 {
// 	// 	fmt.Println("Yellow team has", len(gs.Yellow_team), "4 required")
// 	// 	return ERROR
// 	// } else if len(gs.Blue_team) != 3 {
// 	// 	fmt.Println("Blue team has", len(gs.Blue_team), "3 required")
// 	// 	return ERROR
// 	// }
//
// 	// This assumes that the robot ids range from 0 to m.team_size
// 	for id, at_state := range m.at_states {
// 		target := m.instructionSet[id][at_state].Position
// 		robot := gs.GetRobot(info.ID(id), m.team)
// 		if atPosition(robot, target) {
// 			nr_states := len(m.instructionSet[id])
// 			m.at_states[id] = (at_state + 1) % nr_states // Cycle through the states
// 		}
// 	}
// 	return RUNNING
// }
//
// func atPosition(robot *info.Robot, position info.Position) bool {
//
// 	robot_pos := robot.GetPosition()
// 	diff0 := position.Sub(&robot_pos)
// 	if diff0.Norm() < 100 {
// 		return true
// 	}
// 	return false
// }
//
// // --------------------------MoveToTest----------------------------------------
// type MoveToTest struct {
// 	team     info.Team
// 	at_state int
// 	start    time.Time
// 	max_time time.Duration
// }
//
// func NewMoveToTest(team info.Team) *MoveToTest {
// 	return &MoveToTest{
// 		team:     team,
// 		max_time: 20 * time.Second,
// 		at_state: -1,
// 	}
// }
//
// func (m *MoveToTest) Run() []*info.Instruction {
// 	if m.at_state == -1 {
// 		m.start = time.Now()
// 		m.at_state = 0
// 	}
//
// 	if m.at_state == 0 {
// 		return []*info.Instruction{
// 			{Type: info.MoveToPosition, Id: 0, Position: info.Position{X: 1000, Y: 1000}},
// 			{Type: info.MoveToPosition, Id: 1, Position: info.Position{X: -1000, Y: -1000}},
// 		}
// 	} else if m.at_state == 1 {
// 		return []*info.Instruction{
// 			{Type: info.MoveToPosition, Id: 0, Position: info.Position{X: 100, Y: 100}},
// 			{Type: info.MoveToPosition, Id: 1, Position: info.Position{X: -100, Y: -100}},
// 		}
// 	} else {
// 		return []*info.Instruction{
// 			{Type: info.MoveToPosition, Id: 0, Position: info.Position{X: -1000, Y: -1000}},
// 			{Type: info.MoveToPosition, Id: 1, Position: info.Position{X: 1000, Y: 1000}},
// 		}
// 	}
// }
//
// func (m *MoveToTest) Archived(gs *info.GameState) int {
// 	robot0_pos := gs.GetRobot(info.ID(0), m.team).GetPosition()
// 	robot1_pos := gs.GetRobot(info.ID(1), m.team).GetPosition()
//
// 	if m.at_state == 0 {
// 		target0 := info.Position{X: 1000, Y: 1000}
// 		target1 := info.Position{X: -1000, Y: -1000}
// 		diff0 := target0.Sub(&robot0_pos)
// 		diff1 := target1.Sub(&robot1_pos)
// 		if diff0.Norm() < 100 && diff1.Norm() < 100 {
// 			m.at_state = 1
// 		}
// 	} else if m.at_state == 1 {
// 		target0 := info.Position{X: 100, Y: 100}
// 		target1 := info.Position{X: -100, Y: -100}
// 		diff0 := target0.Sub(&robot0_pos)
// 		diff1 := target1.Sub(&robot1_pos)
// 		if diff0.Norm() < 100 && diff1.Norm() < 100 {
// 			// fmt.Println("norms are", diff0.Norm(), diff1.Norm())
// 			m.at_state = 2
// 		}
// 	} else if m.at_state == 2 {
// 		target0 := info.Position{X: -1000, Y: -1000}
// 		target1 := info.Position{X: 1000, Y: 1000}
// 		diff0 := target0.Sub(&robot0_pos)
// 		diff1 := target1.Sub(&robot1_pos)
// 		if diff0.Norm() < 100 && diff1.Norm() < 100 {
// 			m.at_state = 3
// 		}
// 	}
// 	if m.at_state == 3 {
// 		return COMPLETE
// 	}
// 	if m.at_state >= 0 {
// 		fmt.Println("Time expired", time.Since(m.start))
// 		if time.Since(m.start) > m.max_time {
// 			return TIME_EXPIRED
// 		}
// 	}
//
// 	return RUNNING
// }
//
// // -------------- Move with ball to position -------------------------------------//
//
// type MoveWithBallToPositionTest struct {
// 	team     info.Team
// 	at_state int
// 	start    time.Time
// 	max_time time.Duration
// }
//
// func NewMoveWithBallToPositionTest(team info.Team) *MoveWithBallToPositionTest {
// 	return &MoveWithBallToPositionTest{
// 		team:     team,
// 		max_time: 30 * time.Second,
// 		at_state: -1,
// 	}
// }
//
// func (m *MoveWithBallToPositionTest) Run() []*info.Instruction {
// 	if m.at_state == -1 {
// 		m.start = time.Now()
// 		m.at_state = 0
// 	}
//
// 	if m.at_state == 0 {
// 		return []*info.Instruction{
// 			{Type: info.MoveToBall, Id: 0},
// 		}
// 	} else if m.at_state == 1 {
// 		return []*info.Instruction{
// 			{Type: info.MoveWithBallToPosition, Id: 0, Position: info.Position{X: -2000, Y: -2000}},
// 			{Type: info.MoveToBall, Id: 1},
// 		}
// 	} else {
// 		return []*info.Instruction{
// 			{Type: info.MoveWithBallToPosition, Id: 0, Position: info.Position{X: 0, Y: 0}},
// 			{Type: info.MoveToBall, Id: 1},
// 		}
// 	}
//
// }
//
// func (m *MoveWithBallToPositionTest) Archived(gs *info.GameState) int {
// 	robot_pos := gs.GetRobot(info.ID(0), m.team).GetPosition()
// 	ball_pos := gs.GetBall().GetPosition()
//
// 	dxBall := float64(robot_pos.X - ball_pos.X)
// 	dyBall := float64(robot_pos.Y - ball_pos.Y)
// 	distanceBall := math.Sqrt(math.Pow(dxBall, 2) + math.Pow(dyBall, 2))
//
// 	dxPos := float64(robot_pos.X - (-2000))
// 	dyPos := float64(robot_pos.Y - (-2000))
// 	distancePos := math.Sqrt(math.Pow(dxPos, 2) + math.Pow(dyPos, 2))
//
// 	if m.at_state == 0 {
// 		if distanceBall < 100 {
// 			m.at_state = 1
// 		}
// 	} else if m.at_state == 1 {
// 		if distancePos < 500 {
// 			m.at_state = 2
// 		}
// 	}
//
// 	return RUNNING
// }
//
// // -------------- Kick ball to player -------------------------------------//
//
// type KickToPlayerTest struct {
// 	team     info.Team
// 	at_state int
// 	start    time.Time
// 	max_time time.Duration
// }
//
// func NewKickToPlayerTest(team info.Team) *KickToPlayerTest {
// 	return &KickToPlayerTest{
// 		team:     team,
// 		max_time: 30 * time.Second,
// 		at_state: -1,
// 	}
// }
//
// func (m *KickToPlayerTest) Run() []*info.Instruction {
// 	return []*info.Instruction{}
// 	// if m.at_state == -1 {
// 	// 	m.start = time.Now()
// 	// 	m.at_state = 0
// 	// }
// 	//
// 	// if m.at_state == 0 {
// 	// 	return []*info.Instruction{
// 	// 		{Type: info.MoveToBall, Id: 0},
// 	// 		//{Type: info.MoveToPosition, Id: 2, Position: info.Position{X: -2250, Y: 250}},
// 	// 		{Type: info.MoveToPosition, Id: 1, Position: info.Position{X: -1500, Y: 1500}},
// 	// 		//{Type: info.MoveToPosition, Id: 1, Position: info.Position{X: 2000, Y: 1500}},
// 	// 		//{Type: info.MoveToPosition, Id: 1, Position: info.Position{X: -3000, Y: -2000}},
// 	// 	}
// 	// } else if m.at_state == 1 {
// 	// 	return []*info.Instruction{
// 	// 		{Type: info.MoveWithBallToPosition, Id: 0, Position: info.Position{X: -3000, Y: -1000}},
// 	// 	}
// 	// } else if m.at_state == 2 {
// 	// 	return []*info.Instruction{
// 	// 		{Type: info.KickToPlayer, Id: 0, OtherId: 1},
// 	// 		{Type: info.ReceiveBallFromPlayer, Id: 1, OtherId: 0},
// 	// 	}
// 	// } else {
// 	// 	return []*info.Instruction{
// 	// 		{Type: info.MoveToPosition, Id: 0, Position: info.Position{X: 0, Y: 0}},
// 	// 		{Type: info.ReceiveBallFromPlayer, Id: 1, OtherId: 0},
// 	// 	}
// 	// }
// }
//
// func (m *KickToPlayerTest) Archived(gs *info.GameState) int {
// 	robot_pos0 := gs.GetRobot(info.ID(0), m.team).GetPosition()
// 	robot_pos1 := gs.GetRobot(info.ID(1), m.team).GetPosition()
// 	ball_pos := gs.GetBall().GetPosition()
//
// 	dxBall0 := float64(robot_pos0.X - ball_pos.X)
// 	dyBall0 := float64(robot_pos0.Y - ball_pos.Y)
// 	distanceBall0 := math.Sqrt(math.Pow(dxBall0, 2) + math.Pow(dyBall0, 2))
//
// 	dxPos0 := float64(robot_pos0.X - (-3000))
// 	dyPos0 := float64(robot_pos0.Y - (-1000))
// 	distancePos0 := math.Sqrt(math.Pow(dxPos0, 2) + math.Pow(dyPos0, 2))
//
// 	// The position of the reciver
// 	dxPos1 := float64(robot_pos1.X - (-1500))
// 	dyPos1 := float64(robot_pos1.Y - (1500))
// 	distancePos1 := math.Sqrt(math.Pow(dxPos1, 2) + math.Pow(dyPos1, 2))
//
// 	dxBall1 := float64(robot_pos1.X - ball_pos.X)
// 	dyBall1 := float64(robot_pos1.Y - ball_pos.Y)
// 	distanceBall1 := math.Sqrt(math.Pow(dxBall1, 2) + math.Pow(dyBall1, 2))
//
// 	if m.at_state == 0 {
// 		if distanceBall0 < 100 {
// 			m.at_state = 1
// 		}
// 	} else if m.at_state == 1 {
// 		if distancePos0 < 100 && distancePos1 < 100 {
// 			m.at_state = 2
// 		}
// 	} else if m.at_state == 2 {
// 		if distanceBall1 < 500 {
// 			m.at_state = 3
// 		}
// 	}
//
// 	return RUNNING
// }
