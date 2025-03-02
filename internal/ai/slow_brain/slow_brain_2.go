package ai

import (
	"fmt"
	"sync"
	"time"

	ai "github.com/LiU-SeeGoals/controller/internal/ai/activity"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

// ========================================================
// SlowBrain2 is a demo slow brain with referee handling
// ========================================================

type SlowBrain2 struct {
	SlowBrainComposition
	HandleReferee

	at_state int
	start    time.Time
	max_time time.Duration
	team     info.Team
	prev_ref info.RefCommand
}

func NewSlowBrain2(team info.Team) *SlowBrain2 {
	return &SlowBrain2{
		SlowBrainComposition: SlowBrainComposition{
			team: team,
		},
		HandleReferee: HandleReferee{
			team: team,
		},
		team: team,
	}
}

func (m *SlowBrain2) Init(
	incoming <-chan info.GameInfo,
	activities *[]ai.Activity,
	lock *sync.Mutex,
	team info.Team,
) {
	m.incomingGameInfo = incoming
	m.activities = activities // store pointer directly
	m.activity_lock = lock
	m.start = time.Now()

	go m.run()
}

// This is the main loop of the AI in this slow brain
func (m *SlowBrain2) run() {
	way_points := []info.Position{
		{X: 0, Y: 0, Z: 0, Angle: 0},
		{X: 0, Y: 1000, Z: 0, Angle: 0},
		{X: 1000, Y: 0, Z: 0, Angle: 0},
	}
	index := 0

	for {
		// No need for slow brain to be fast
		time.Sleep(100 * time.Millisecond)

		gameInfo := <-m.incomingGameInfo

		referee_activities := m.GetRefereeActivities(&gameInfo)
		fmt.Println("referee action: ", gameInfo.Status.GetGameEvent().RefCommand)
		if referee_activities != nil {
			m.ReplaceActivities(referee_activities)
			m.at_state = REFEREE
			continue
		}

		if m.at_state == REFEREE {
			m.ClearActivities()
			m.at_state = RUNNING
		}

		if len(*m.activities) == 0 {
			fmt.Println("done with action: ", m.team)
			m.AddActivity(ai.NewMoveToPosition(m.team, 0, way_points[index]))
			index = (index + 1) % len(way_points)
		}
	}
}
