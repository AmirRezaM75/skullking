package ws

// TODO: Constant for all card ids

type Factory struct {
}

func (f Factory) skullKing() Card {
	return Card{
		Id:     1,
		Color:  "#000",
		Number: 0,
		Group:  "skullKing",
	}
}

func (f Factory) whale() Card {
	return Card{
		Id:     2,
		Color:  "#fff",
		Number: 0,
		Group:  "whale",
	}
}

func (f Factory) kraken() Card {
	return Card{
		Id:     3,
		Color:  "#ff1760",
		Number: 0,
		Group:  "kraken",
	}
}

func (f Factory) mermaids() []Card {
	var cards []Card

	id := 4

	for i := 1; i <= 2; i++ {
		card := Card{
			Id:     id,
			Color:  "#fffded",
			Number: i,
			Group:  "mermaid",
		}
		cards = append(cards, card)
		id++
	}

	return cards
}

func (f Factory) parrots() []Card {
	var cards []Card

	id := 10

	for i := 1; i <= 14; i++ {
		card := Card{
			Id:     id,
			Color:  "#69ff69",
			Number: i,
			Group:  "parrot",
		}
		cards = append(cards, card)
		id++
	}

	return cards
}

// 14 Pirate Maps
func (f Factory) maps() []Card {
	var cards []Card

	id := 30

	for i := 1; i <= 14; i++ {
		card := Card{
			Id:     id,
			Color:  "#6e5eff",
			Number: i,
			Group:  "map",
		}
		cards = append(cards, card)
		id++
	}

	return cards
}

// 14 Treasure Chests
func (f Factory) chests() []Card {
	var cards []Card

	id := 50

	for i := 1; i <= 14; i++ {
		card := Card{
			Id:     id,
			Color:  "#ffeb57",
			Number: i,
			Group:  "chest",
		}
		cards = append(cards, card)
		id++
	}

	return cards
}

// 14 Jolly Roger cards
func (f Factory) roger() []Card {
	var cards []Card

	id := 70

	for i := 1; i <= 14; i++ {
		card := Card{
			Id:     id,
			Color:  "#000",
			Number: i,
			Group:  "roger",
		}
		cards = append(cards, card)
		id++
	}

	return cards
}

func (f Factory) pirates() []Card {
	var cards []Card

	id := 90

	for i := 1; i <= 5; i++ {
		card := Card{
			Id:     id,
			Color:  "#b53564",
			Number: i,
			Group:  "pirate",
		}
		cards = append(cards, card)
		id++
	}

	return cards
}

func (f Factory) escape() []Card {
	var cards []Card

	id := 100

	for i := 1; i <= 5; i++ {
		card := Card{
			Id:     id,
			Color:  "#fffded",
			Number: i,
			Group:  "mermaid",
		}
		cards = append(cards, card)
		id++
	}

	return cards
}
