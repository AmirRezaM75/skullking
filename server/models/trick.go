package models

type Trick struct {
	Number int
	// This is useful to find out who is eligible when receiving 'PICK' command.
	PickingUserId string
	// We cannot use a map for picked cards because the order in which the cards are picked
	// affects the pickable cards. In Go, a map is not a suitable collection to maintain
	// the sequence of elements as it is unordered.
	PickedCards        []PickedCard
	WinnerPlayerId     string // TODO: Shorter (Winner)
	StarterPlayerIndex int
}

type PickedCard struct {
	PlayerId string
	CardId   CardId
}

func (trick Trick) isPlayerPicked(playerId string) bool {
	for _, pickedCard := range trick.PickedCards {
		if playerId == pickedCard.PlayerId {
			return true
		}
	}

	return false
}

func (trick Trick) getPickedCardByPlayerId(playerId string) *PickedCard {
	for _, pickedCard := range trick.PickedCards {
		if playerId == pickedCard.PlayerId {
			return &pickedCard
		}
	}
	return nil
}
