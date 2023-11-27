package datatypes

// Used for sending actions to grsim
type Parameters struct {
	// Yellow team?
	YellowTeam bool
	// Which robot?
	RobotId uint32
	// Robot's forward speed
	VelTangent float32
	// Robot's side speed
	VelNormal float32
	// Robot's rotational speed
	VelAngular float32
	// Kick
	KickSpeedX float32
	// DONT USE
	KickSpeedZ float32
	// Dribbler
	Spinner bool
	//DONT USE
	WheelsSpeed bool
}

// Creates and defaults parameters to 0
func NewParameters() *Parameters {

	return &Parameters{YellowTeam: true,
		VelTangent:  0,
		VelNormal:   0,
		VelAngular:  0,
		KickSpeedX:  0,
		KickSpeedZ:  0,
		Spinner:     false,
		WheelsSpeed: false,
	}
}
