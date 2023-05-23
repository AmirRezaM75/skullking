package models

type Table struct {
	cards []Card
}

func (t Table) suit() Card {
	var suit Card

	for _, card := range t.cards {
		if card.isSuit() {
			suit = card
		}
	}

	return suit
}
