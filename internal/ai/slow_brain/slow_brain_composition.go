package ai

import (
	"fmt"
	"sync"
	"time"

	ai "github.com/LiU-SeeGoals/controller/internal/ai/activity"
	"github.com/LiU-SeeGoals/controller/internal/info"
	. "github.com/LiU-SeeGoals/controller/internal/logger"
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
	activity_lock    *sync.Mutex                  // shared mutex for synchronization
	prev_ref         bool
	waiting_for_kick bool
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

func (m *SlowBrainComposition) ReplaceActivities(activities *[info.TEAM_SIZE]ai.Activity) {
	m.activity_lock.Lock()
	defer m.activity_lock.Unlock()
	m.activities = activities
}

func (m *SlowBrainComposition) HandleRef(gi *info.GameInfo, active_robots []int) bool {

	fmt.Println(gi.Status.GetGameEvent().GetCurrentState())

	switch gi.Status.GetGameEvent().GetCurrentState() {
	case info.STATE_HALTED, info.STATE_STOPPED, info.STATE_PENALTY_PREPARATION, info.STATE_FREE_KICK, info.STATE_TIMEOUT, info.STATE_BALL_PLACEMENT:
		for _, value := range active_robots {
			m.AddActivity(ai.NewStop(info.ID(value)))
		}
		m.prev_ref = true
		return true

	case info.STATE_KICKOFF_PREPARATION:
		if gi.Status.GetGameEvent().GetTeamWithPossession() == m.team { // We are kicker
			m.AddActivity(ai.NewRefKicker(info.ID(1), m.team))
			// m.AddActivity(ai.NewRefKickoff(info.ID(0), m.team))
		} else { // We are not kicker
			m.AddActivity(ai.NewRefKickoff(info.ID(1), m.team))
		}
		m.waiting_for_kick = true
		m.prev_ref = true
		return true

	default:
		if m.waiting_for_kick {
			if gi.State.Ball.GetVelocity().Norm() > 10 {
				m.waiting_for_kick = false
				m.prev_ref = false
				return false
			} else {
				m.prev_ref = true
				return true
			}
		}
		// If we are exiting ref activity
		if m.prev_ref == true {
			m.ClearActivities()
		}
		m.prev_ref = false
		return false
	}
}
