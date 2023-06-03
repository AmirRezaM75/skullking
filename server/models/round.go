package models

type Round struct {
	Number             int
	DealtCards         map[string][]CardId
	Bids               map[string]int
	Tricks             map[int]*Trick // TODO: Simple slice?
	StarterPlayerIndex int
}

func (round Round) getDealtCardIdsByPlayerId(playerId string) []int {
	var cardIds []int

	for _, cardId := range round.DealtCards[playerId] {
		cardIds = append(cardIds, int(cardId))
	}

	return cardIds
}

func (round Round) getPickedCardIdsByPlayerId(playerId string) []CardId {
	var cardIds []CardId

	for _, trick := range round.Tricks {
		for _, pickedCard := range trick.PickedCards {
			if playerId == pickedCard.PlayerId {
				cardIds = append(cardIds, pickedCard.CardId)
			}
		}
	}

	return cardIds
}
