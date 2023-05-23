package models

type Round struct {
	Number         int
	Scores         map[string]int
	DealtCards     map[string][]CardId
	RemainingCards map[string][]CardId
	Bids           map[string]int
	Tricks         []Trick
}
