package models

type Table struct {
	cards []Card
}

func newTable(cardIds []CardId) Table {
	var cards []Card

	for _, cardId := range cardIds {
		cards = append(cards, newCardFromId(cardId))
	}

	return Table{
		cards: cards,
	}
}

func (t Table) pattern() string {
	var pattern string

	if len(t.cards) == 0 {
		return ""
	}

	if t.cards[0].isCharacter() || t.cards[0].isBeast() {
		return ""
	}

	for _, card := range t.cards {
		if card.isSuit() {
			pattern = card.Type
		}
	}

	return pattern
}
