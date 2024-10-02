package state

type Zone struct {
	scores []float32
}
type GameAnalisys struct {
	team      state.Team
	myTeam    [][]Zone // 2D array of zones
	otherTeam [][]Zone // 2D array of zones
	zoneSize  float32
}
