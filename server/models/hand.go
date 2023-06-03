package models

type Hand struct {
	cards []Card
}

func newHand(cardIds []CardId) Hand {
	var cards []Card

	var card Card

	for _, cardId := range cardIds {
		cards = append(cards, card.fromId(cardId))
	}

	return Hand{
		cards: cards,
	}
}

func (s Hand) pickables(t Table) []CardId {
	var specialIds []CardId

	var cardIds []CardId

	var options []CardId

	pattern := t.pattern()

	for _, card := range s.cards {

		cardIds = append(cardIds, card.Id)

		if card.isSpecial() {
			specialIds = append(specialIds, card.Id)
		}

		if card.Type == pattern {
			options = append(options, card.Id)
		}
	}

	if len(options) == 0 {
		return cardIds
	}

	return append(options, specialIds...)
}
