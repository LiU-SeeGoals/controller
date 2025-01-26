package ai

import (
	"sync"
	"time"

	ai "github.com/LiU-SeeGoals/controller/internal/ai/activity"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type SlowBrain1 struct {
	SlowBrainComposition
	team          info.Team
	at_state      int
	start         time.Time
	max_time      time.Duration
	activity      []ai.Activity
	activity_lock *sync.Mutex // Shared mutex for synchronization
}

func NewSlowBrain1(team info.Team) *SlowBrain1 {
	return &SlowBrain1{
		team: team,
	}
}

func (m *SlowBrain1) Init(incoming <-chan info.GameInfo,
	activities *[]ai.Activity,
	lock *sync.Mutex,
	team info.Team,
) {
	m.incomingGameInfo = incoming
	m.activity = *activities
	m.activity_lock = lock
	m.team = team
	m.start = time.Now()
	go m.run()
}

// This is the main loop of the AI in this slow brain
func (m *SlowBrain1) run() {
	way_points := []info.Position{
		info.Position{X: 0, Y: 0, Z: 0, Angle: 0},
		info.Position{X: 0, Y: 1000, Z: 0, Angle: 0},
		info.Position{X: 1000, Y: 0, Z: 0, Angle: 0},
	}
	index := 0

	for {
		if len(m.activity) == 0 {
			m.AddActivity(ai.NewMoveToPosition(m.team, 2, way_points[index]))
			index = (index + 1) % len(way_points)
		}
	}

}
