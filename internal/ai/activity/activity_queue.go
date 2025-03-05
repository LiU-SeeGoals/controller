package ai

import (
	"fmt"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type ActivityQueue struct {
	activities []Activity
	current    int
	id        info.ID
}

func (q *ActivityQueue) String() string {
	return fmt.Sprintf("ActivityQueue: %s", q.activities[q.current].String())
}

func NewActivityQueue(id info.ID, activities []Activity) *ActivityQueue {
	return &ActivityQueue{
		activities: activities,
		current:    0,
	}
}

func (q *ActivityQueue) GetAction(gi *info.GameInfo) action.Action {
	if q.activities[q.current].Achieved(gi) {
		q.current += 1
	}

	return q.activities[q.current].GetAction(gi)
}

func (q *ActivityQueue) Achieved(gi *info.GameInfo) bool {
	return q.current == len(q.activities) - 1 && q.activities[q.current].Achieved(gi)
}

func (q *ActivityQueue) GetID() info.ID {
	return q.id
}
