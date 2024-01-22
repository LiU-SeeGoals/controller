package ai

type PlayFinder struct {
}

func NewPlayFinder() *PlayFinder {
	pf := &PlayFinder{}
	return pf
}

type Plays struct {
}

func (pf *PlayFinder) FindPlays(gameAnalysis *GameAnalysis) *Plays {

	return &Plays{}
}
