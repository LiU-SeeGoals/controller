package ai

import "fmt"

type PlayFinder struct {
}

func NewPlayFinder() *PlayFinder {
	pf := &PlayFinder{}
	return pf
}

type Plays struct {
}

func (pf *PlayFinder) FindPlays(data *Data) *Plays {
	fmt.Println("Strutsen")
	return &Plays{}
}
