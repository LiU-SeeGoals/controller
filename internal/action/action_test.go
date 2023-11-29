package action

import (
	"testing"

	"github.com/LiU-SeeGoals/controller/internal/proto/basestation"
)

// Not done. Need to add more tests for the other actions.

func TestStopTranslateReal(t *testing.T) {
	stop := Stop{Id: 42}

	// Call the TranslateReal method
	command := stop.TranslateReal()

	// Expected result
	expectedCommand := &basestation.Command{
		CommandId: basestation.ActionType_STOP_ACTION,
		RobotId:   42,
	}

	if (command.GetCommandId() != expectedCommand.GetCommandId()) || (command.GetRobotId() != expectedCommand.GetRobotId()) {
		t.Errorf("TranslateReal result is not as expected")
	}
}
