package models

// Set Owned user's cards
type Set struct {
	cards []Card
}

func (s Set) pickables(t Table) []CardId {
	var specialIds []CardId

	var cardIds []CardId

	var options []CardId

	var suit Card

	suit = t.suit()

	for _, card := range s.cards {

		cardIds = append(cardIds, card.Id)

		if card.isSpecial() {
			specialIds = append(specialIds, card.Id)
		}

		if card.Type == suit.Type {
			options = append(options, card.Id)
		}
	}

	if len(options) == 0 {
		return cardIds
	}

	return append(options, specialIds...)
}
