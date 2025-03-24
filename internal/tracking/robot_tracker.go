package tracking

import (
	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/info"
)


type robotTracker struct {
	clientSource  *client.TrackerSource
	latestFrameNum uint32
	team 	info.Team
	id   	info.ID
}

func NewRobotTracker(source *client.TrackerSource, team info.Team, id info.ID) *robotTracker {
	return &robotTracker{clientSource: source, team: team, id: id}
}

func (t *robotTracker) GetTrackedRobot() (info.RobotState) {
	robotState, frameNum, ok := t.clientSource.GetTrackedRobot(t.team, t.id)
	
	if !ok || frameNum == t.latestFrameNum {
		// No new data
		return info.RobotState{Valid: false}
	}
	t.latestFrameNum = frameNum
	return robotState
}

