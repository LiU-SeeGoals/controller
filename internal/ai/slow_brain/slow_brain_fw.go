package ai

import (
	"fmt"
	"sync"
	"time"
	ai "github.com/LiU-SeeGoals/controller/internal/ai/activity"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type SlowBrainFw struct {
	SlowBrainComposition
	at_state int
	start    time.Time
	max_time time.Duration
}

func NewSlowBrainFw(team info.Team) *SlowBrain1 {
	return &SlowBrain1{
		SlowBrainComposition: SlowBrainComposition{
			team: team,
		},
	}
}

func (m *SlowBrainFw) Init(
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

// This is the main loop of the AI in this slow brain
func (m *SlowBrainFw) run() {
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
		{X: -2920, Y: -4100, Z: 0, Angle: 0},
		{X: -5900, Y: -1950, Z: 0, Angle: 0},
		{X: -4250, Y: -1950, Z: 0, Angle: 0},
	}
	index := 0
    robot := 0

	for {
        //sgameInfo := <-m.incomingGameInfo
        //sfmt.Println(gameInfo.Status)

		// No need for slow brain to be fast
		time.Sleep(100 * time.Millisecond)

		if m.activities[robot] == nil {
			fmt.Println("done with action: ", m.team)
			fmt.Println("next action: ", way_points[index])
			m.AddActivity(ai.NewMoveToPosition(m.team, info.ID(robot), way_points[index]))
			index = (index + 1) % len(way_points)
		}
	}
}
