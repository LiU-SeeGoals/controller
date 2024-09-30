package state

type GameAnalisys struct {
	Blue_team   [TEAM_SIZE]*Robot
	Yellow_team [TEAM_SIZE]*Robot

	// Holds ball data
	Ball *Ball

	MessageReceived int64
	LagTime         int64
}
