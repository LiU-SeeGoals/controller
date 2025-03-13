package ai

import (
	"fmt"
	"sync"
	"time"

	ai "github.com/LiU-SeeGoals/controller/internal/ai/activity"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type SlowBrain1 struct {
	SlowBrainComposition
	at_state int
	start    time.Time
	max_time time.Duration
}

func NewSlowBrain1(team info.Team) *SlowBrain1 {
	return &SlowBrain1{
		SlowBrainComposition: SlowBrainComposition{
			team: team,
		},
	}
}

func (m *SlowBrain1) Init(
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
func (m *SlowBrain1) run() {

	way_points := []info.Position{
        // Go between line
		//{X: -3575, Y: -4128, Z: 0, Angle: 0},
		//{X: -5558, Y: -4096, Z: 0, Angle: 0},
        // Go to pos
		//{X: -4195, Y: -3544, Z: 0, Angle: 0},
        // Triangle
		{X: -5500, Y: -4100, Z: 0, Angle: 0},
		{X: -5600, Y: -2600, Z: 0, Angle: 0},
		{X: -4200, Y: -3400, Z: 0, Angle: 0},
	}
	index := 0

	for {
		// No need for slow brain to be fast
		time.Sleep(100 * time.Millisecond)

		//fmt.Println("slow, number of activities:", len(*m.activities))

		if m.activities[1] == nil {
			fmt.Println("done with action: ", m.team)
            //fmt.Println("Next move action: ", way_points[index])
			//time.Sleep(10000 * time.Millisecond)
			//m.AddActivity(ai.NewMoveToPosition(m.team, 1, way_points[index]))
            action := ai.NewMoveToBall(m.team, 1)
            //fmt.Println("Move to ball: ", action.GetAction())
			m.AddActivity(action)
			index = (index + 1) % len(way_points)
		}
	}
}
