package tracking

import (
	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

const (
	visibilityThreshold = 0.5
)

type ballTracker struct {
	clientSource  *client.TrackerSource
	latestFrameNum uint32
}

func NewBallTracker(source *client.TrackerSource) *ballTracker {
	return &ballTracker{clientSource: source}
}

func (t *ballTracker) GetTrackedBall(lastPossessor *info.Robot) (info.BallState, bool) {
	trackedBall, frameNum, ok := t.clientSource.GetTrackedBall()

	if frameNum == t.latestFrameNum || !ok {
		// No new data
		return trackedBall, false
	}
	t.latestFrameNum = frameNum

	// If ball is occluded and we have a possessor, track the ball with the possessor
	if trackedBall.Visibility < visibilityThreshold && lastPossessor != nil {
		ball := info.NewBallState(
			lastPossessor.DribblerPos(),
			lastPossessor.GetVelocity(),
			trackedBall.Visibility,
			trackedBall.Timestamp,
			"occluded",
		)
		return ball, ok
	}

	return trackedBall, ok
}

