package models

import (
	"math/rand"
	"time"
)

// Deck A complete pack, or deck, includes 14 cards in each suit + special cards
type Deck struct {
	Cards []Card
}

func (d Deck) Shuffle() {
	cards := d.Cards

	for i := range cards {
		source := rand.NewSource(time.Now().UnixNano())
		r := rand.New(source)
		j := r.Intn(len(d.Cards) - 1)
		cards[i], cards[j] = cards[j], cards[i]
	}
}

func (d Deck) Deal(count, size int) [][]CardId {
	var output [][]CardId

	index := 0

	for i := 0; i < count; i++ {
		var cardIds []CardId

		cards := d.Cards[index : size+index]

		for _, card := range cards {
			cardIds = append(cardIds, card.Id)
		}

		output = append(output, cardIds)

		index = size + index
	}

	return output
}
