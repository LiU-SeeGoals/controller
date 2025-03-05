package ai

import (
	"fmt"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type ActivityLoop struct {
	activities []Activity
	current    int
	id        info.ID
}

func (l *ActivityLoop) String() string {
	return fmt.Sprintf("ActivityLoop: %s", l.activities[l.current].String())
}

func NewActivityLoop(id info.ID, activities []Activity) *ActivityLoop {
	return &ActivityLoop{
		activities: activities,
		current:    0,
	}
}

func (l *ActivityLoop) GetAction(gi *info.GameInfo) action.Action {
	if l.activities[l.current].Achieved(gi) {
		l.current = (l.current + 1) % len(l.activities)
	}

	return l.activities[l.current].GetAction(gi)
}

func (l *ActivityLoop) Achieved(gi *info.GameInfo) bool {
	return false
}

func (l *ActivityLoop) GetID() info.ID {
	return l.id
}

