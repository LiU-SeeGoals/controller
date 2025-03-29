package ai

import (
	"fmt"
	"sync"
	"time"

	ai "github.com/LiU-SeeGoals/controller/internal/ai/activity"
	"github.com/LiU-SeeGoals/controller/internal/info"
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

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
 *                                                                                       *
 *                                                                                       *
 * This is Rasmus Wallin's file, touch it without asking and you shall meet your demise! *
 *                                                                                       *
 *                                                                                       *
 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

func (m *SlowBrainAo) run() {
	way_points := []info.Position{
		// Go between line
		//{X: -3575, Y: -4128, Z: 0, Angle: 0},
		//{X: -5558, Y: -4096, Z: 0, Angle: 0},
		// Go to pos
		//{X: -4195, Y: -3544, Z: 0, Angle: 0},
		// Triangle
		//{X: -5500, Y: -4100, Z: 0, Angle: 0},
		//{X: -5600, Y: -2600, Z: 0, Angle: 0},
		//{X: -4200, Y: -3400, Z: 0, Angle: 0},
		// Triangle 2
		{X: -2000, Y: 0, Z: 0, Angle: 0},
		{X: 500, Y: 2000, Z: 0, Angle: 150},
		{X: 2000, Y: 1000, Z: 0, Angle: 100},
	}
	index := 0
	succesfull_commands := 0
	robots := []int{0, 1, 3}

	gameInfo := <-m.incomingGameInfo
	fmt.Println(gameInfo.Status)
	// Basic idea
	// Defender: Get some dudes to guard the goal, stand in "line" formation towards ball
	//	- Function that returns indices for robots that should perform defense
	// Attacker: Chase ball, kick toward goal, turn to support when away from ball;
	// Support: Stand a bit away from attack so he can pass, turn into attacker when get ball

	for {
		// No need for slow brain to be fast
		time.Sleep(100 * time.Millisecond)

		//if m.HandleRef(&gameInfo, robots) {
		//	continue
		//}

        robot := robots[0]
        if m.activities[robot] == nil {

			posx := float64(gameInfo.Field.GetFieldLength()/2 - gameInfo.Field.GetGoalWidth())

			pos := info.Position{X: posx, Y: 0, Z: 0, Angle: 3.14}

			ball, err := gameInfo.State.Ball.GetEstimatedPosition()

			if err != nil {
				fmt.Println("failed ball")
				ball = pos
			}

            fmt.Println(fmt.Sprintf("done with (%d) action (%s)", robot, m.team))
            fmt.Println("next action: ", way_points[index])
            fmt.Println("sent commands: ", succesfull_commands)
            m.AddActivity(ai.NewMoveToPosition(m.team, info.ID(robot), ball))
            index = (index + 1) % len(way_points)
            succesfull_commands++
        }

        robot = robots[1]
        if m.activities[robot] == nil {
            fmt.Println(fmt.Sprintf("done with (%d) action (%s)", robot, m.team))
            fmt.Println("next action: ", way_points[index])
            fmt.Println("sent commands: ", succesfull_commands)
            m.AddActivity(ai.NewMoveToPosition(m.team, info.ID(robot), way_points[index]))
            index = (index + 1) % len(way_points)
            succesfull_commands++
        }

        robot = robots[2]
        if m.activities[robot] == nil {
            fmt.Println(fmt.Sprintf("done with (%d) action (%s)", robot, m.team))
            fmt.Println("next action: ", way_points[index])
            fmt.Println("sent commands: ", succesfull_commands)
            m.AddActivity(ai.NewMoveToPosition(m.team, info.ID(robot), way_points[index]))
            index = (index + 1) % len(way_points)
            succesfull_commands++
        }
	}
}
