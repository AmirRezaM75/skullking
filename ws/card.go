package ws

type Card struct {
	// TODO: Change to CardId
	Id     CardId `json:"id"`
	Color  string `json:"color"`
	Number int    `json:"number"`
	Group  string `json:"group"`
}

type CardId int

// It looks stupid at first to define each single of them as const
// But it makes it easier to understand what's the card in winner() func for example.
const (
	SkullKing CardId = iota + 1
	Whale
	Kraken
	Mermaid1
	Mermaid2
	Parrot1
	Parrot2
	Parrot3
	Parrot4
	Parrot5
	Parrot6
	Parrot7
	Parrot8
	Parrot9
	Parrot10
	Parrot11
	Parrot12
	Parrot13
	Parrot14
	Map1
	Map2
	Map3
	Map4
	Map5
	Map6
	Map7
	Map8
	Map9
	Map10
	Map11
	Map12
	Map13
	Map14
	Chest1
	Chest2
	Chest3
	Chest4
	Chest5
	Chest6
	Chest7
	Chest8
	Chest9
	Chest10
	Chest11
	Chest12
	Chest13
	Chest14
	Roger1
	Roger2
	Roger3
	Roger4
	Roger5
	Roger6
	Roger7
	Roger8
	Roger9
	Roger10
	Roger11
	Roger12
	Roger13
	Roger14
	Pirate1
	Pirate2
	Pirate3
	Pirate4
	Pirate5
	Escape1
	Escape2
	Escape3
	Escape4
	Escape5
)

// I found it more performant to predefine cards before starting game,
// Instead of generating them whenever I need to filter...
var cards = map[CardId]Card{
	SkullKing: {
		Id:     SkullKing,
		Color:  "#000",
		Number: 0,
		Group:  "king",
	},
	Whale: {
		Id:     Whale,
		Color:  "#fff",
		Number: 0,
		Group:  "whale",
	},
	Kraken: {
		Id:     Kraken,
		Color:  "#ff1760",
		Number: 0,
		Group:  "kraken",
	},
	Mermaid1: {
		Id:     Mermaid1,
		Color:  "#fffded",
		Number: 0,
		Group:  "mermaid",
	},
	Mermaid2: {
		Id:     Mermaid2,
		Color:  "#fffded",
		Number: 0,
		Group:  "mermaid",
	},
	// 14 Parrots
	Parrot1: {
		Id:     Parrot1,
		Color:  "#69ff69",
		Number: 1,
		Group:  "parrot",
	},
	Parrot2: {
		Id:     Parrot2,
		Color:  "#69ff69",
		Number: 2,
		Group:  "parrot",
	},
	Parrot3: {
		Id:     Parrot3,
		Color:  "#69ff69",
		Number: 3,
		Group:  "parrot",
	},
	Parrot4: {
		Id:     Parrot4,
		Color:  "#69ff69",
		Number: 4,
		Group:  "parrot",
	},
	Parrot5: {
		Id:     Parrot5,
		Color:  "#69ff69",
		Number: 5,
		Group:  "parrot",
	},
	Parrot6: {
		Id:     Parrot6,
		Color:  "#69ff69",
		Number: 6,
		Group:  "parrot",
	},
	Parrot7: {
		Id:     Parrot7,
		Color:  "#69ff69",
		Number: 7,
		Group:  "parrot",
	},
	Parrot8: {
		Id:     Parrot8,
		Color:  "#69ff69",
		Number: 8,
		Group:  "parrot",
	},
	Parrot9: {
		Id:     Parrot9,
		Color:  "#69ff69",
		Number: 9,
		Group:  "parrot",
	},
	Parrot10: {
		Id:     Parrot10,
		Color:  "#69ff69",
		Number: 10,
		Group:  "parrot",
	},
	Parrot11: {
		Id:     Parrot11,
		Color:  "#69ff69",
		Number: 11,
		Group:  "parrot",
	},
	Parrot12: {
		Id:     Parrot12,
		Color:  "#69ff69",
		Number: 12,
		Group:  "parrot",
	},
	Parrot13: {
		Id:     Parrot13,
		Color:  "#69ff69",
		Number: 13,
		Group:  "parrot",
	},
	Parrot14: {
		Id:     Parrot14,
		Color:  "#69ff69",
		Number: 14,
		Group:  "parrot",
	},
	// 14 Pirate Maps
	Map1: {
		Id:     Map1,
		Color:  "#6e5eff",
		Number: 1,
		Group:  "map",
	},
	Map2: {
		Id:     Map2,
		Color:  "#6e5eff",
		Number: 2,
		Group:  "map",
	},
	Map3: {
		Id:     Map3,
		Color:  "#6e5eff",
		Number: 3,
		Group:  "map",
	},
	Map4: {
		Id:     Map4,
		Color:  "#6e5eff",
		Number: 4,
		Group:  "map",
	},
	Map5: {
		Id:     Map5,
		Color:  "#6e5eff",
		Number: 5,
		Group:  "map",
	},
	Map6: {
		Id:     Map6,
		Color:  "#6e5eff",
		Number: 6,
		Group:  "map",
	},
	Map7: {
		Id:     Map7,
		Color:  "#6e5eff",
		Number: 7,
		Group:  "map",
	},
	Map8: {
		Id:     Map8,
		Color:  "#6e5eff",
		Number: 8,
		Group:  "map",
	},
	Map9: {
		Id:     Map9,
		Color:  "#6e5eff",
		Number: 9,
		Group:  "map",
	},
	Map10: {
		Id:     Map10,
		Color:  "#6e5eff",
		Number: 10,
		Group:  "map",
	},
	Map11: {
		Id:     Map11,
		Color:  "#6e5eff",
		Number: 11,
		Group:  "map",
	},
	Map12: {
		Id:     Map12,
		Color:  "#6e5eff",
		Number: 12,
		Group:  "map",
	},
	Map13: {
		Id:     Map13,
		Color:  "#6e5eff",
		Number: 13,
		Group:  "map",
	},
	Map14: {
		Id:     Map14,
		Color:  "#6e5eff",
		Number: 14,
		Group:  "map",
	},
	// 14 Treasure Chests
	Chest1: {
		Id:     Chest1,
		Color:  "#ffeb57",
		Number: 1,
		Group:  "chest",
	},
	Chest2: {
		Id:     Chest2,
		Color:  "#ffeb57",
		Number: 2,
		Group:  "chest",
	},
	Chest3: {
		Id:     Chest3,
		Color:  "#ffeb57",
		Number: 3,
		Group:  "chest",
	},
	Chest4: {
		Id:     Chest4,
		Color:  "#ffeb57",
		Number: 4,
		Group:  "chest",
	},
	Chest5: {
		Id:     Chest5,
		Color:  "#ffeb57",
		Number: 5,
		Group:  "chest",
	},
	Chest6: {
		Id:     Chest6,
		Color:  "#ffeb57",
		Number: 6,
		Group:  "chest",
	},
	Chest7: {
		Id:     Chest7,
		Color:  "#ffeb57",
		Number: 7,
		Group:  "chest",
	},
	Chest8: {
		Id:     Chest8,
		Color:  "#ffeb57",
		Number: 8,
		Group:  "chest",
	},
	Chest9: {
		Id:     Chest9,
		Color:  "#ffeb57",
		Number: 9,
		Group:  "chest",
	},
	Chest10: {
		Id:     Chest10,
		Color:  "#ffeb57",
		Number: 10,
		Group:  "chest",
	},
	Chest11: {
		Id:     Chest11,
		Color:  "#ffeb57",
		Number: 11,
		Group:  "chest",
	},
	Chest12: {
		Id:     Chest12,
		Color:  "#ffeb57",
		Number: 12,
		Group:  "chest",
	},
	Chest13: {
		Id:     Chest13,
		Color:  "#ffeb57",
		Number: 13,
		Group:  "chest",
	},
	Chest14: {
		Id:     Chest14,
		Color:  "#ffeb57",
		Number: 14,
		Group:  "chest",
	},
	// 14 Jolly Roger cards
	Roger1: {
		Id:     Roger1,
		Color:  "#000",
		Number: 1,
		Group:  "roger",
	},
	Roger2: {
		Id:     Roger2,
		Color:  "#000",
		Number: 2,
		Group:  "roger",
	},
	Roger3: {
		Id:     Roger3,
		Color:  "#000",
		Number: 3,
		Group:  "roger",
	},
	Roger4: {
		Id:     Roger4,
		Color:  "#000",
		Number: 4,
		Group:  "roger",
	},
	Roger5: {
		Id:     Roger5,
		Color:  "#000",
		Number: 5,
		Group:  "roger",
	},
	Roger6: {
		Id:     Roger6,
		Color:  "#000",
		Number: 6,
		Group:  "roger",
	},
	Roger7: {
		Id:     Roger7,
		Color:  "#000",
		Number: 7,
		Group:  "roger",
	},
	Roger8: {
		Id:     Roger8,
		Color:  "#000",
		Number: 8,
		Group:  "roger",
	},
	Roger9: {
		Id:     Roger9,
		Color:  "#000",
		Number: 9,
		Group:  "roger",
	},
	Roger10: {
		Id:     Roger10,
		Color:  "#000",
		Number: 10,
		Group:  "roger",
	},
	Roger11: {
		Id:     Roger11,
		Color:  "#000",
		Number: 11,
		Group:  "roger",
	},
	Roger12: {
		Id:     Roger12,
		Color:  "#000",
		Number: 12,
		Group:  "roger",
	},
	Roger13: {
		Id:     Roger13,
		Color:  "#000",
		Number: 13,
		Group:  "roger",
	},
	Roger14: {
		Id:     Roger14,
		Color:  "#000",
		Number: 14,
		Group:  "roger",
	},
	// 5 Pirates
	Pirate1: {
		Id:     Pirate1,
		Color:  "#b53564",
		Number: 0,
		Group:  "pirate",
	},
	Pirate2: {
		Id:     Pirate2,
		Color:  "#b53564",
		Number: 0,
		Group:  "pirate",
	},
	Pirate3: {
		Id:     Pirate3,
		Color:  "#b53564",
		Number: 0,
		Group:  "pirate",
	},
	Pirate4: {
		Id:     Pirate4,
		Color:  "#b53564",
		Number: 0,
		Group:  "pirate",
	},
	Pirate5: {
		Id:     Pirate5,
		Color:  "#b53564",
		Number: 0,
		Group:  "pirate",
	},
	// 5 Escape cards
	Escape1: {
		Id:     Escape1,
		Color:  "#fffded",
		Number: 0,
		Group:  "escape",
	},
	Escape2: {
		Id:     Escape2,
		Color:  "#fffded",
		Number: 0,
		Group:  "escape",
	},
	Escape3: {
		Id:     Escape3,
		Color:  "#fffded",
		Number: 0,
		Group:  "escape",
	},
	Escape4: {
		Id:     Escape4,
		Color:  "#fffded",
		Number: 0,
		Group:  "escape",
	},
	Escape5: {
		Id:     Escape5,
		Color:  "#fffded",
		Number: 0,
		Group:  "escape",
	},
}

func winner(cardIds []CardId) CardId {
	var lead Card

	// Suit cards are the numbered cards, 1-14, in four colors.
	var suitCardLead Card

	for _, id := range cardIds {
		card := cards[id]

		if card.Group == "parrot" || card.Group == "map" || card.Group == "chest" || card.Group == "roger" {
			if suitCardLead.Id == 0 || suitCardLead.Number < card.Number {
				suitCardLead = card
			}
		}

		if lead.Id == 0 {
			lead = card
			continue
		}

		if lead.Group == card.Group {
			if lead.Number < card.Number {
				lead = card
			}
		} else {

			if card.Group == "whale" || card.Group == "kraken" {
				lead = card
			}

			if (card.Group == "parrot" || card.Group == "map" || card.Group == "chest") &&
				lead.Group == "escape" {
				lead = card
			}

			if (card.Group == "roger" || card.Group == "mermaid" || card.Group == "pirate" || card.Group == "king") &&
				(lead.Group == "parrot" || lead.Group == "map" || lead.Group == "chest" || lead.Group == "escape") {
				lead = card
			}

			if card.Group == "pirate" &&
				(lead.Group == "roger" || lead.Group == "mermaid") {
				lead = card
			}

			if card.Group == "king" && lead.Group == "pirate" {
				lead = card
			}

			if card.Group == "mermaid" && lead.Group == "king" {
				lead = card
			}
		}
	}

	if lead.Group == "whale" {
		return suitCardLead.Id
	}

	return lead.Id
}
