package datatypes

type Parameters struct {
	YellowTeam  bool
	RobotId     uint32
	VelTangent  float32
	VelNormal   float32
	VelAngular  float32
	KickSpeedX  float32
	KickSpeedZ  float32
	Spinner     bool
	WheelsSpeed bool
}

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
