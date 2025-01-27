package ai

import (
	"fmt"
	"sync"
	"time"

	"github.com/LiU-SeeGoals/controller/internal/action"
	ai "github.com/LiU-SeeGoals/controller/internal/ai/activity"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type FastBrainGO struct {
	team             info.Team
	incomingGameInfo <-chan info.GameInfo
	outgoingActions  chan<- []action.Action
	activities       *[]ai.Activity // <-- pointer to a slice
	activity_lock    *sync.Mutex    // shared mutex for synchronization
}

func NewFastBrainGO() *FastBrainGO {
	return &FastBrainGO{}
}

func (fb *FastBrainGO) Init(
	incoming <-chan info.GameInfo,
	activities *[]ai.Activity,
	lock *sync.Mutex,
	outgoing chan<- []action.Action,
	team info.Team,
) {
	fb.incomingGameInfo = incoming
	fb.outgoingActions = outgoing
	fb.team = team
	fb.activity_lock = lock

	// Store the pointer directly
	fb.activities = activities

	go fb.Run()
}

func (fb *FastBrainGO) Run() {
	for {
		// For example, throttle the loop slightly to avoid busy-loop:
		time.Sleep(50 * time.Millisecond) // or read from fb.incomingGameInfo if event-driven

		gameInfo := <-fb.incomingGameInfo
		// Make a snapshot of current activities under lock
		fb.activity_lock.Lock()
		activitiesCopy := make([]ai.Activity, len(*fb.activities))
		copy(activitiesCopy, *fb.activities)
		fb.activity_lock.Unlock()

		var actions []action.Action
		for i := range activitiesCopy {
			// If done, remove it from the *shared* slice
			if activitiesCopy[i].Achieved(&gameInfo) {
				fmt.Println("sucessful action")
				fb.activity_lock.Lock()
				// find it in the real slice (not in the copy!)
				for j, realAct := range *fb.activities {
					if realAct == activitiesCopy[i] {
						*fb.activities = append(
							(*fb.activities)[:j],
							(*fb.activities)[j+1:]...,
						)
						break
					}
				}
				fb.activity_lock.Unlock()
			} else {
				// Otherwise, get an action
				actions = append(actions, activitiesCopy[i].GetAction(&gameInfo))
			}
		}

		// Send actions
		fb.outgoingActions <- actions
	}
}
