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
	activities *[info.TEAM_SIZE]ai.Activity,
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
	fmt.Println("slow brain 2: starting")

	for {
		// No need for slow brain to be fast
		time.Sleep(100 * time.Millisecond)

		gameInfo := <-m.incomingGameInfo
		fmt.Println("slow brain 2: ")

		referee_activities := m.GetRefereeActivities(&gameInfo)
		fmt.Println("referee action: ", gameInfo.Status.GetGameEvent().RefCommand)

		for _, act := range referee_activities {
			if act != nil { // Check for a non-nil activity
				fmt.Println("adding referee activity")
				m.ReplaceActivities(referee_activities)
				m.at_state = REFEREE
			}
		}

		// If we are EXITING the REFEREE state, we need to clear the activities
		if m.at_state == REFEREE {
			fmt.Println("clearing activities")
			m.ClearActivities()
			m.at_state = RUNNING
		}

		if m.activities[0] == nil {
			fmt.Println("done with action: ", m.team)
			m.AddActivity(ai.NewMoveToPosition(m.team, 0, way_points[index]))
			index = (index + 1) % len(way_points)
		}
	}
}
