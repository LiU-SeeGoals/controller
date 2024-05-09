package ai

type StrategyFinder struct {
}

func NewPlayFinder() *StrategyFinder {
	pf := &StrategyFinder{}
	return pf
}

type Plays struct {
}

func (pf *StrategyFinder) FindStrategy(gameAnalysis *GameAnalysis) *Plays {

	return &Plays{}
}
