package info

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

type RefCommand int

const (
	// All robots should completely stop moving.
	HALT RefCommand = iota
	// Robots must keep 50 cm from the ball.
	STOP
	// A prepared kickoff or penalty may now be taken.
	NORMAL_START
	// The ball is dropped and free for either team.
	FORCE_START
	// The yellow team may move into kickoff position.
	PREPARE_KICKOFF_YELLOW
	// The blue team may move into kickoff position.
	PREPARE_KICKOFF_BLUE
	// The yellow team may move into penalty position.
	PREPARE_PENALTY_YELLOW
	// The blue team may move into penalty position.
	PREPARE_PENALTY_BLUE
	// The yellow team may take a direct free kick.
	DIRECT_FREE_YELLOW
	// The blue team may take a direct free kick.
	DIRECT_FREE_BLUE
	// The yellow team may take an indirect free kick.
	INDIRECT_FREE_YELLOW
	// The blue team may take an indirect free kick.
	INDIRECT_FREE_BLUE
	// The yellow team is currently in a timeout.
	TIMEOUT_YELLOW
	// The blue team is currently in a timeout.
	TIMEOUT_BLUE
	GOAL_YELLOW // DEPRICATED
	GOAL_BLUE // DEPRICATED
	// Equivalent to STOP, but the yellow team must pick up the ball and
	// drop it in the Designated Position.
	BALL_PLACEMENT_YELLOW
	// Equivalent to STOP, but the blue team must pick up the ball and drop
	// it in the Designated Position.
	BALL_PLACEMENT_BLUE
)

// // MatchEvent is an enum that represents the current match event.
// type MatchState int
//
// const (
// 	// The match has been halted, robots must stop moving.
// 	HALTED MatchState = iota
// 	// The match is stopped, robots may move
// 	STOPPED
//
// 	// The match is preparing for a kickoff by the yellow team.
// 	PREPARING_KICKOFF_YELLOW
// 	// The match is preparing for a kickoff by the blue team.
// 	PREPARING_KICKOFF_BLUE
// 	// Yellow team is making a kickoff
// 	KICKOFF_YELLOW
// 	// Blue team is making a kickoff
// 	KICKOFF_BLUE
//
// 	// The match is preparing for a penalty kick by the yellow team.
// 	PREPARING_PENALTY_YELLOW
// 	// The match is preparing for a penalty kick by the blue team.
// 	PREPARING_PENALTY_BLUE
// 	// Yellow team is making a penalty kick
// 	PENALTY_KICK_YELLOW
// 	// Blue team is making a penalty kick
// 	PENALTY_KICK_BLUE
//
// 	// Yellow team is placing ball
// 	PLACING_BALL_YELLOW
// 	// Blue team is placing ball
// 	PLACING_BALL_BLUE
// 	// Yellow team is taking a free kick
// 	TAKING_FREE_KICK_YELLOW
// 	// Blue team is taking a free kick
// 	TAKING_FREE_KICK_BLUE
//
// 	// Yellow team is taking a timeout.
// 	TAKING_TIMEOUT_YELLOW
// 	// Blue team is taking a timeout.
// 	TAKING_TIMEOUT_BLUE
//
// 	// The match is running.
// 	RUNNING
// )

type GameEvent struct {


	// Command issued by the referee.
	RefCommand RefCommand
	// The UNIX timestamp when the command was issued, in microseconds.
	// This value changes only when a new command is issued, not on each packet.
	Command_timestamp uint64
	// The coordinates of the Designated Position. These are measured in
	// millimetres and correspond to SSL-Vision coordinates. These are
	// present (in the case of a ball placement command) or
	// absent (in the case of any other command).
	DesignatedPosition *mat.VecDense
	// The command that will be issued after the current stoppage and ball placement to continue the game.
	next_command RefCommand
	// The time in microseconds that is remaining until the current action times out
	// The time will not be reset. It can get negative.
	// An autoRef would raise an appropriate event, if the time gets negative.
	// Possible actions where this time is relevant:
	//  * free kicks
	//  * kickoff, penalty kick, force start
	//  * ball placement
	current_action_time_remaining int64

	// All game events that were detected since the last RUNNING info.
	// Will be cleared as soon as the game is continued.
	// game_events GameEvent
	// All proposed game events that were detected since the last RUNNING info.
	// game_event_proposals GameEventProposalGroup

}

func NewGameEvent() *GameEvent {
	return &GameEvent{
		DesignatedPosition: mat.NewVecDense(2, nil),
	}
}


// String method for RefCommand to convert the enum to a human-readable string
func (rc RefCommand) String() string {
	switch rc {
	case HALT:
		return "Halt"
	case STOP:
		return "Stop"
	case NORMAL_START:
		return "Normal Start"
	case FORCE_START:
		return "Force Start"
	case PREPARE_KICKOFF_YELLOW:
		return "Prepare Kickoff Yellow"
	case PREPARE_KICKOFF_BLUE:
		return "Prepare Kickoff Blue"
	case PREPARE_PENALTY_YELLOW:
		return "Prepare Penalty Yellow"
	case PREPARE_PENALTY_BLUE:
		return "Prepare Penalty Blue"
	case DIRECT_FREE_YELLOW:
		return "Direct Free Yellow"
	case DIRECT_FREE_BLUE:
		return "Direct Free Blue"
	case INDIRECT_FREE_YELLOW:
		return "Indirect Free Yellow"
	case INDIRECT_FREE_BLUE:
		return "Indirect Free Blue"
	case TIMEOUT_YELLOW:
		return "Timeout Yellow"
	case TIMEOUT_BLUE:
		return "Timeout Blue"
	case BALL_PLACEMENT_YELLOW:
		return "Ball Placement Yellow"
	case BALL_PLACEMENT_BLUE:
		return "Ball Placement Blue"
	default:
		return fmt.Sprintf("Unknown RefCommand (%d)", rc)
	}
}

// String method for GameEvent
func (ge *GameEvent) String() string {
	// Format DesignatedPosition if it's not nil
	position := "N/A"
	if ge.DesignatedPosition != nil {
		position = fmt.Sprintf("(x: %.2f, y: %.2f)", ge.DesignatedPosition.At(0, 0), ge.DesignatedPosition.At(1, 0))
	}

	// Create a formatted string for the GameEvent
	return fmt.Sprintf(
		"Game Event:\n"+
			"  Ref Command: %s\n"+
			"  Command Timestamp: %d microseconds\n"+
			"  Designated Position: %s\n"+
			"  Next Command: %s\n"+
			"  Current Action Time Remaining: %d microseconds",
		ge.RefCommand.String(),
		ge.Command_timestamp,
		position,
		ge.next_command.String(),
		ge.current_action_time_remaining,
	)
}

// Getter and Setter for RefCommand
func (ge *GameEvent) GetRefCommand() RefCommand {
	return ge.RefCommand
}

func (ge *GameEvent) SetRefCommand(command RefCommand) {
	ge.RefCommand = command
}

// Getter and Setter for CommandTimestamp
func (ge *GameEvent) GetCommandTimestamp() uint64 {
	return ge.Command_timestamp
}

func (ge *GameEvent) SetCommandTimestamp(timestamp uint64) {
	ge.Command_timestamp = timestamp
}

// Getter and Setter for DesignatedPosition
func (ge *GameEvent) GetDesignatedPosition() *mat.VecDense {
	return ge.DesignatedPosition
}

func (ge *GameEvent) SetDesignatedPosition(x float64, y float64) {
	ge.DesignatedPosition.SetVec(0, x)
	ge.DesignatedPosition.SetVec(1, y)
}

// Getter and Setter for NextCommand
func (ge *GameEvent) GetNextCommand() RefCommand {
	return ge.next_command
}

func (ge *GameEvent) SetNextCommand(command RefCommand) {
	ge.next_command = command
}

// Getter and Setter for CurrentActionTimeRemaining
func (ge *GameEvent) GetCurrentActionTimeRemaining() int64 {
	return ge.current_action_time_remaining
}

func (ge *GameEvent) SetCurrentActionTimeRemaining(timeRemaining int64) {
	ge.current_action_time_remaining = timeRemaining
}
