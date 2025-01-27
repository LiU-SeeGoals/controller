package ai

import (
	"time"

	ai "github.com/LiU-SeeGoals/controller/internal/ai/activity"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

const (
	RUNNING int = iota
	COMPLETE
	TIME_EXPIRED
	ERROR
	FAILED
)

type SlowBrainComposition struct {
	team             info.Team
	incomingGameInfo <-chan info.GameInfo
	scale            float32
	run_scenario     int // -1 for all
	start            time.Time
}

func (m *SlowBrain1) ClearActivities() {
	m.activity_lock.Lock()
	defer m.activity_lock.Unlock()
	*m.activities = []ai.Activity{}
}

func (m *SlowBrain1) AddActivity(activity ai.Activity) {
	m.activity_lock.Lock()
	defer m.activity_lock.Unlock()
	*m.activities = append(*m.activities, activity)
}
