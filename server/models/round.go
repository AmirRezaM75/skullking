package models

type Round struct {
	Number             int
	Scores             map[string]int
	DealtCards         map[string][]CardId
	RemainingCards     map[string][]CardId
	Bids               map[string]int
	Tricks             map[int]*Trick
	StarterPlayerIndex int
}

func (round Round) getDealtCardIdsByPlayerId(playerId string) []int {
	var cardIds []int

	for _, cardId := range round.DealtCards[playerId] {
		cardIds = append(cardIds, int(cardId))
	}

	return cardIds
}
