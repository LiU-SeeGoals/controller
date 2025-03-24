package receiver

import (
	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/config"
	"github.com/LiU-SeeGoals/controller/internal/info"
	"github.com/LiU-SeeGoals/controller/internal/tracking"
	"github.com/LiU-SeeGoals/proto_go/gc"
	"github.com/LiU-SeeGoals/proto_go/ssl_vision"
)

type SSLReceiver struct {
	vision   *client.SSLVisionSource
	referee  *client.SSLRefereeSource
	tracking *tracking.Tracking
}

func NewSSLReceiver() *SSLReceiver {
	return &SSLReceiver{
		vision:   client.NewSSLVisionSource(config.GetSSLClientAddress()),
		referee:  client.NewSSLRefereeSource(config.GetGCClientAddress()),
		tracking: tracking.NewTracking(config.GetSSLTrackerAddress()),
	}
}

func (receiver *SSLReceiver) Update(gameInfo *info.GameInfo, playTime int64) {
	// detection := receiver.vision.GetVisionData()

	receiver.tracking.UpdateTracking(gameInfo.State)

	refereeData := receiver.referee.GetRefereeData()
	UpdateReferee(gameInfo.Status, refereeData)
}

func UpdateReferee(status *info.GameStatus, refereeData *gc.Referee) {
	if refereeData == nil {
		return
	}

	status.SetGameEvent(
		info.RefCommand(refereeData.GetCommand().Number()),
		refereeData.GetCommandTimestamp(),
		float64(refereeData.GetDesignatedPosition().GetX()),
		float64(refereeData.GetDesignatedPosition().GetY()),
		info.RefCommand(refereeData.GetNextCommand().Number()),
		refereeData.GetCurrentActionTimeRemaining())

	status.SetGameStatus(info.GameStage(refereeData.GetStage().Number()),
		info.MatchType(refereeData.GetMatchType().Number()),
		refereeData.GetPacketTimestamp(),
		refereeData.GetStageTimeLeft(),
		refereeData.GetCommandCounter(),
		refereeData.GetBlueTeamOnPositiveHalf(),
		refereeData.GetStatusMessage())

	// yellow team
	yellow := refereeData.GetYellow()
	if yellow != nil {
		status.SetTeamInfo(
			true,
			yellow.GetName(),
			yellow.GetScore(),
			yellow.GetRedCards(),
			yellow.GetYellowCards(),
			yellow.GetTimeouts(),
			yellow.GetTimeoutTime(),
			yellow.GetGoalkeeper(),
			yellow.GetFoulCounter(),
			yellow.GetBallPlacementFailures(),
			yellow.GetMaxAllowedBots(),
			yellow.GetBotSubstitutionsLeft(),
			yellow.GetBotSubstitutionTimeLeft(),
			yellow.GetYellowCardTimes(),
			yellow.GetCanPlaceBall(),
			yellow.GetBotSubstitutionIntent(),
			yellow.GetBallPlacementFailuresReached(),
			yellow.GetBotSubstitutionAllowed(),
		)
	}

	// blue team
	blue := refereeData.GetBlue()

	if blue != nil {
		status.SetTeamInfo(
			false,
			blue.GetName(),
			blue.GetScore(),
			blue.GetRedCards(),
			blue.GetYellowCards(),
			blue.GetTimeouts(),
			blue.GetTimeoutTime(),
			blue.GetGoalkeeper(),
			blue.GetFoulCounter(),
			blue.GetBallPlacementFailures(),
			blue.GetMaxAllowedBots(),
			blue.GetBotSubstitutionsLeft(),
			blue.GetBotSubstitutionTimeLeft(),
			blue.GetYellowCardTimes(),
			blue.GetCanPlaceBall(),
			blue.GetBotSubstitutionIntent(),
			blue.GetBallPlacementFailuresReached(),
			blue.GetBotSubstitutionAllowed(),
		)
	}
}

// Update the geometry of the field
func (receiver *SSLReceiver) UpdateGeometry(field *info.GameField, geometry *ssl_vision.SSL_GeometryData) {
	if geometry == nil {
		return
	}

	fieldData := geometry.GetField()

	field.SetField(fieldData.GetFieldLength(),
		fieldData.GetFieldWidth(),
		fieldData.GetGoalWidth(),
		fieldData.GetGoalDepth(),
		fieldData.GetBoundaryWidth(),
		fieldData.GetPenaltyAreaDepth(),
		fieldData.GetPenaltyAreaWidth(),
	)

	if field.FieldLines == nil {
		for _, line := range fieldData.GetFieldLines() {
			field.AddFieldLine(line.GetName(), float64(line.GetP1().GetX()), float64(line.GetP1().GetY()), float64(line.GetP2().GetX()), float64(line.GetP2().GetY()), float64(line.GetThickness()), int(line.GetType()))
		}
		for _, arc := range fieldData.GetFieldArcs() {
			field.AddFieldArc(arc.GetName(), float64(arc.GetCenter().GetX()), float64(arc.GetCenter().GetY()), float64(arc.GetRadius()), float64(arc.GetA1()), float64(arc.GetA2()), float64(arc.GetThickness()), int(arc.GetType()))
		}
	}
}
