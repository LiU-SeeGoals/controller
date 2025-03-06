package ai

import (
	"sync"
	"time"

	. "github.com/LiU-SeeGoals/controller/internal/logger"
	ai "github.com/LiU-SeeGoals/controller/internal/ai/activity"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

const (
	RUNNING int = iota
	COMPLETE
	TIME_EXPIRED
	ERROR
	FAILED
	REFEREE
)

type SlowBrainComposition struct {
	team             info.Team
	incomingGameInfo <-chan info.GameInfo
	scale            float64
	run_scenario     int // -1 for all
	start            time.Time
	activities       *[info.TEAM_SIZE]ai.Activity // <-- pointer to the slice
	activity_lock    *sync.Mutex    // shared mutex for synchronization
}

func (m *SlowBrainComposition) ClearActivities() {
	m.activity_lock.Lock()
	defer m.activity_lock.Unlock()
	*m.activities = [info.TEAM_SIZE]ai.Activity{}
}

func (m *SlowBrainComposition) AddActivity(activity ai.Activity) {
	m.activity_lock.Lock()
	defer m.activity_lock.Unlock()
	idx := activity.GetID()
	Logger.Infof("Adding activity %v", activity)
	m.activities[idx] = activity
}

func (m *SlowBrainComposition) ReplaceActivities(activities [info.TEAM_SIZE]ai.Activity) {
	m.activity_lock.Lock()
	defer m.activity_lock.Unlock()
	m.activities = &activities
}
