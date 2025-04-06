package ai

import (
	"fmt"
	"math"
	"sync"
	"time"

	ai "github.com/LiU-SeeGoals/controller/internal/ai/activity"
	"github.com/LiU-SeeGoals/controller/internal/info"
	"github.com/LiU-SeeGoals/controller/internal/logger"
	vis "github.com/LiU-SeeGoals/controller/internal/visualisation"
	"gonum.org/v1/plot/plotter"
)

type SlowBrainAo struct {
	SlowBrainComposition
	at_state int
	start    time.Time
	max_time time.Duration
}

func NewSlowBrainAo(team info.Team) *SlowBrainAo {
	return &SlowBrainAo{
		SlowBrainComposition: SlowBrainComposition{
			team: team,
		},
	}
}

func (m *SlowBrainAo) Init(
	incoming <-chan info.GameInfo,
	activities *[info.TEAM_SIZE]ai.Activity,
	lock *sync.Mutex,
	team info.Team,
) {
	m.incomingGameInfo = incoming
	m.activities = activities // store pointer directly
	m.activity_lock = lock
	m.team = team
	m.start = time.Now()

	go m.run()
}

func (m *SlowBrainAo) run() {

	gameInfo := <-m.incomingGameInfo
	fmt.Println(gameInfo.Status)

	// Basic idea
	// Defender: Get some dudes to guard the goal, stand in "line" formation towards ball
	//	- Function that returns indices for robots that should perform defense
	// Attacker: Chase ball, kick toward goal, turn to support when away from ball;
	// Support: Stand a bit away from attack so he can pass, turn into attacker when get ball
	robotPos := plotter.XYs{}

	fig := vis.GetVisualiser().CreateEmptyPlotWindow()
	for {
		// No need for slow brain to be fast
		time.Sleep(1 * time.Millisecond)

		robots := []int{0,1,3}
		if m.HandleRef(&gameInfo, robots) {
			continue
		}

        // robot := robots[0]
		defenders := []info.ID{0,1}
		attackers := []info.ID{3}

		myRobotPos, err := gameInfo.State.GetTeam(m.team)[0].GetPosition()
		if err != nil {
			logger.Logger.Debugln("Big err")
		}

		robotPos = append(robotPos, plotter.XY{X: myRobotPos.X, Y: myRobotPos.Y})
		p := vis.ScatterPlt(robotPos)
		p.Title.Text = fmt.Sprintf("Robot %v team %v", 0, m.team)
		fig.UpdatePlotWindow(p)

		m.defense(defenders)
		m.attack(attackers)

			//      if m.activities[robot] == nil {
			//
			//          fmt.Println(fmt.Sprintf("done with (%d) action (%s)", robot, m.team))
			//          fmt.Println("next action: ", way_points[index])
			//          fmt.Println("sent commands: ", succesfull_commands)
			//
			// activityLoop := []ai.Activity{
			// 	ai.NewMoveToBall(m.team, 0),
			// 	ai.NewKickAtPosition(m.team, 0, info.Position{X: 2000, Y: 2000, Z: 0, Angle: 0}),
			// 	// ai.NewKickToPlayer(m.team, 0, 1),
			// }
			// loop := ai.NewActivityLoop(0, activityLoop)
			// m.AddActivity(loop)
			//
			//          index = (index + 1) % len(way_points)
			//          succesfull_commands++
			//      }
	}
}

func (m *SlowBrainAo) defense(robots []info.ID){

	gi := <-m.incomingGameInfo

	var formation = map[info.ID][2]float64{
		robots[0]: {0, 0},
		robots[1]: {0, -200},
		// robots[2]: {0, 200},
	}

	def := gi.HomeGoalLine(m.team)
	ballpos, err := gi.State.GetBall().GetPosition()
	defY := ballpos.Y

	if err != nil {
		fmt.Println("Ball position is undefined")
	}

	defensePos := info.Position{X: def.X, Y: defY, Z: 0, Angle: def.Angle + math.Pi}

	for i := range robots {
		id := robots[i]
		offset := formation[id]
		// fmt.Printf("robots %v i %v id %v offest %v", robots, i, id, offset)
		formationPosx := defensePos.X + offset[0]
		formationPosy := defensePos.Y + offset[1]
		pos := info.Position{X: formationPosx, Y: formationPosy, Z: 0, Angle: defensePos.Angle}
		// fmt.Printf("Moving %v to %v\n", id, pos)
		m.AddActivity(ai.NewMoveToPosition(m.team, id, pos))
	}
}

func (m *SlowBrainAo) attack(robots []info.ID){

	for i := range robots{
		// if m.activities[robots[i]] == nil {
			activityLoop := []ai.Activity{
				ai.NewMoveToBall(m.team, robots[i]),
				ai.NewKickTheBall(m.team, robots[i], info.Position{X: 2000, Y: 2000, Z: 0, Angle: 0}),
				// ai.NewKickToPlayer(m.team, 0, 1),
			}
			loop := ai.NewActivityLoop(robots[i], activityLoop)
			m.AddActivity(loop)
		// }
	}
}
