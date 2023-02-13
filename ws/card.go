package ws

type Card struct {
	Id     CardId `json:"id"`
	Color  string `json:"color"`
	Number int    `json:"number"`
	Type   string `json:"type"`
}

func (c Card) fromId(id CardId) Card {
	return cards[id]
}

// Suit cards are the numbered cards, 1-14, in four colors.
func (c Card) isSuit() bool {
	return c.isStandardSuit() || c.isRoger()
}

// There are three standard suits; Parrot (green), Treasure Chest (yellow), Treasure Map (purple),
// and the trump suit: Jolly Roger (Black)
func (c Card) isStandardSuit() bool {
	return c.isParrot() || c.isMap() || c.isChest()
}

func (c Card) isKing() bool {
	return c.Type == TypeKing
}

func (c Card) isWhale() bool {
	return c.Type == TypeWhale
}

func (c Card) isKraken() bool {
	return c.Type == TypeKraken
}

func (c Card) isMermaid() bool {
	return c.Type == TypeMermaid
}

func (c Card) isSpecial() bool {
	return c.isKing() || c.isMermaid() || c.isPirate() || c.isEscape() || c.isKraken() || c.isWhale()
}

func (c Card) isParrot() bool {
	return c.Type == TypeParrot
}

func (c Card) isMap() bool {
	return c.Type == TypeMap
}

func (c Card) isChest() bool {
	return c.Type == TypeChest
}

func (c Card) isRoger() bool {
	return c.Type == TypeRoger
}

func (c Card) isPirate() bool {
	return c.Type == TypePirate
}

func (c Card) isEscape() bool {
	return c.Type == TypeEscape
}

const TypeKing string = "king"
const TypeWhale string = "whale"
const TypeKraken string = "kraken"
const TypeMermaid string = "mermaid"
const TypeParrot string = "parrot"
const TypeMap string = "map"
const TypeChest string = "chest"
const TypeRoger string = "roger"
const TypePirate string = "pirate"
const TypeEscape string = "escape"

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

// TODO: Should be part of Deck struct
// I found it more performant to predefine cards before starting game,
// Instead of generating them whenever I need to filter...
var cards = map[CardId]Card{
	SkullKing: {
		Id:     SkullKing,
		Color:  "#000",
		Number: 0,
		Type:   TypeKing,
	},
	Whale: {
		Id:     Whale,
		Color:  "#fff",
		Number: 0,
		Type:   TypeWhale,
	},
	Kraken: {
		Id:     Kraken,
		Color:  "#ff1760",
		Number: 0,
		Type:   TypeKraken,
	},
	Mermaid1: {
		Id:     Mermaid1,
		Color:  "#fffded",
		Number: 0,
		Type:   TypeMermaid,
	},
	Mermaid2: {
		Id:     Mermaid2,
		Color:  "#fffded",
		Number: 0,
		Type:   TypeMermaid,
	},
	// 14 Parrots
	Parrot1: {
		Id:     Parrot1,
		Color:  "#69ff69",
		Number: 1,
		Type:   TypeParrot,
	},
	Parrot2: {
		Id:     Parrot2,
		Color:  "#69ff69",
		Number: 2,
		Type:   TypeParrot,
	},
	Parrot3: {
		Id:     Parrot3,
		Color:  "#69ff69",
		Number: 3,
		Type:   TypeParrot,
	},
	Parrot4: {
		Id:     Parrot4,
		Color:  "#69ff69",
		Number: 4,
		Type:   TypeParrot,
	},
	Parrot5: {
		Id:     Parrot5,
		Color:  "#69ff69",
		Number: 5,
		Type:   TypeParrot,
	},
	Parrot6: {
		Id:     Parrot6,
		Color:  "#69ff69",
		Number: 6,
		Type:   TypeParrot,
	},
	Parrot7: {
		Id:     Parrot7,
		Color:  "#69ff69",
		Number: 7,
		Type:   TypeParrot,
	},
	Parrot8: {
		Id:     Parrot8,
		Color:  "#69ff69",
		Number: 8,
		Type:   TypeParrot,
	},
	Parrot9: {
		Id:     Parrot9,
		Color:  "#69ff69",
		Number: 9,
		Type:   TypeParrot,
	},
	Parrot10: {
		Id:     Parrot10,
		Color:  "#69ff69",
		Number: 10,
		Type:   TypeParrot,
	},
	Parrot11: {
		Id:     Parrot11,
		Color:  "#69ff69",
		Number: 11,
		Type:   TypeParrot,
	},
	Parrot12: {
		Id:     Parrot12,
		Color:  "#69ff69",
		Number: 12,
		Type:   TypeParrot,
	},
	Parrot13: {
		Id:     Parrot13,
		Color:  "#69ff69",
		Number: 13,
		Type:   TypeParrot,
	},
	Parrot14: {
		Id:     Parrot14,
		Color:  "#69ff69",
		Number: 14,
		Type:   TypeParrot,
	},
	// 14 Pirate Maps
	Map1: {
		Id:     Map1,
		Color:  "#6e5eff",
		Number: 1,
		Type:   TypeMap,
	},
	Map2: {
		Id:     Map2,
		Color:  "#6e5eff",
		Number: 2,
		Type:   TypeMap,
	},
	Map3: {
		Id:     Map3,
		Color:  "#6e5eff",
		Number: 3,
		Type:   TypeMap,
	},
	Map4: {
		Id:     Map4,
		Color:  "#6e5eff",
		Number: 4,
		Type:   TypeMap,
	},
	Map5: {
		Id:     Map5,
		Color:  "#6e5eff",
		Number: 5,
		Type:   TypeMap,
	},
	Map6: {
		Id:     Map6,
		Color:  "#6e5eff",
		Number: 6,
		Type:   TypeMap,
	},
	Map7: {
		Id:     Map7,
		Color:  "#6e5eff",
		Number: 7,
		Type:   TypeMap,
	},
	Map8: {
		Id:     Map8,
		Color:  "#6e5eff",
		Number: 8,
		Type:   TypeMap,
	},
	Map9: {
		Id:     Map9,
		Color:  "#6e5eff",
		Number: 9,
		Type:   TypeMap,
	},
	Map10: {
		Id:     Map10,
		Color:  "#6e5eff",
		Number: 10,
		Type:   TypeMap,
	},
	Map11: {
		Id:     Map11,
		Color:  "#6e5eff",
		Number: 11,
		Type:   TypeMap,
	},
	Map12: {
		Id:     Map12,
		Color:  "#6e5eff",
		Number: 12,
		Type:   TypeMap,
	},
	Map13: {
		Id:     Map13,
		Color:  "#6e5eff",
		Number: 13,
		Type:   TypeMap,
	},
	Map14: {
		Id:     Map14,
		Color:  "#6e5eff",
		Number: 14,
		Type:   TypeMap,
	},
	// 14 Treasure Chests
	Chest1: {
		Id:     Chest1,
		Color:  "#ffeb57",
		Number: 1,
		Type:   TypeChest,
	},
	Chest2: {
		Id:     Chest2,
		Color:  "#ffeb57",
		Number: 2,
		Type:   TypeChest,
	},
	Chest3: {
		Id:     Chest3,
		Color:  "#ffeb57",
		Number: 3,
		Type:   TypeChest,
	},
	Chest4: {
		Id:     Chest4,
		Color:  "#ffeb57",
		Number: 4,
		Type:   TypeChest,
	},
	Chest5: {
		Id:     Chest5,
		Color:  "#ffeb57",
		Number: 5,
		Type:   TypeChest,
	},
	Chest6: {
		Id:     Chest6,
		Color:  "#ffeb57",
		Number: 6,
		Type:   TypeChest,
	},
	Chest7: {
		Id:     Chest7,
		Color:  "#ffeb57",
		Number: 7,
		Type:   TypeChest,
	},
	Chest8: {
		Id:     Chest8,
		Color:  "#ffeb57",
		Number: 8,
		Type:   TypeChest,
	},
	Chest9: {
		Id:     Chest9,
		Color:  "#ffeb57",
		Number: 9,
		Type:   TypeChest,
	},
	Chest10: {
		Id:     Chest10,
		Color:  "#ffeb57",
		Number: 10,
		Type:   TypeChest,
	},
	Chest11: {
		Id:     Chest11,
		Color:  "#ffeb57",
		Number: 11,
		Type:   TypeChest,
	},
	Chest12: {
		Id:     Chest12,
		Color:  "#ffeb57",
		Number: 12,
		Type:   TypeChest,
	},
	Chest13: {
		Id:     Chest13,
		Color:  "#ffeb57",
		Number: 13,
		Type:   TypeChest,
	},
	Chest14: {
		Id:     Chest14,
		Color:  "#ffeb57",
		Number: 14,
		Type:   TypeChest,
	},
	// 14 Jolly Roger cards
	Roger1: {
		Id:     Roger1,
		Color:  "#000",
		Number: 1,
		Type:   TypeRoger,
	},
	Roger2: {
		Id:     Roger2,
		Color:  "#000",
		Number: 2,
		Type:   TypeRoger,
	},
	Roger3: {
		Id:     Roger3,
		Color:  "#000",
		Number: 3,
		Type:   TypeRoger,
	},
	Roger4: {
		Id:     Roger4,
		Color:  "#000",
		Number: 4,
		Type:   TypeRoger,
	},
	Roger5: {
		Id:     Roger5,
		Color:  "#000",
		Number: 5,
		Type:   TypeRoger,
	},
	Roger6: {
		Id:     Roger6,
		Color:  "#000",
		Number: 6,
		Type:   TypeRoger,
	},
	Roger7: {
		Id:     Roger7,
		Color:  "#000",
		Number: 7,
		Type:   TypeRoger,
	},
	Roger8: {
		Id:     Roger8,
		Color:  "#000",
		Number: 8,
		Type:   TypeRoger,
	},
	Roger9: {
		Id:     Roger9,
		Color:  "#000",
		Number: 9,
		Type:   TypeRoger,
	},
	Roger10: {
		Id:     Roger10,
		Color:  "#000",
		Number: 10,
		Type:   TypeRoger,
	},
	Roger11: {
		Id:     Roger11,
		Color:  "#000",
		Number: 11,
		Type:   TypeRoger,
	},
	Roger12: {
		Id:     Roger12,
		Color:  "#000",
		Number: 12,
		Type:   TypeRoger,
	},
	Roger13: {
		Id:     Roger13,
		Color:  "#000",
		Number: 13,
		Type:   TypeRoger,
	},
	Roger14: {
		Id:     Roger14,
		Color:  "#000",
		Number: 14,
		Type:   TypeRoger,
	},
	// 5 Pirates
	Pirate1: {
		Id:     Pirate1,
		Color:  "#b53564",
		Number: 0,
		Type:   TypePirate,
	},
	Pirate2: {
		Id:     Pirate2,
		Color:  "#b53564",
		Number: 0,
		Type:   TypePirate,
	},
	Pirate3: {
		Id:     Pirate3,
		Color:  "#b53564",
		Number: 0,
		Type:   TypePirate,
	},
	Pirate4: {
		Id:     Pirate4,
		Color:  "#b53564",
		Number: 0,
		Type:   TypePirate,
	},
	Pirate5: {
		Id:     Pirate5,
		Color:  "#b53564",
		Number: 0,
		Type:   TypePirate,
	},
	// 5 Escape cards
	Escape1: {
		Id:     Escape1,
		Color:  "#fffded",
		Number: 0,
		Type:   TypeEscape,
	},
	Escape2: {
		Id:     Escape2,
		Color:  "#fffded",
		Number: 0,
		Type:   TypeEscape,
	},
	Escape3: {
		Id:     Escape3,
		Color:  "#fffded",
		Number: 0,
		Type:   TypeEscape,
	},
	Escape4: {
		Id:     Escape4,
		Color:  "#fffded",
		Number: 0,
		Type:   TypeEscape,
	},
	Escape5: {
		Id:     Escape5,
		Color:  "#fffded",
		Number: 0,
		Type:   TypeEscape,
	},
}

// TODO: Should have Table struct as receiver
func winner(cardIds []CardId) CardId {
	var lead Card

	var suitLead Card

	var mermaidLead Card

	hasPirate := false

	hasKing := false

	for _, id := range cardIds {
		card := cards[id]

		// Instead of traversing card items more than once
		// We update existence flag of specific cards here.

		if card.isPirate() {
			hasPirate = true
		}

		if card.isKing() {
			hasKing = true
		}

		// Define leaders
		if card.isSuit() {
			if suitLead.Id == 0 || suitLead.Number < card.Number {
				suitLead = card
			}
		}

		if card.isMermaid() && mermaidLead.Id == 0 {
			mermaidLead = card
		}

		if lead.Id == 0 {
			lead = card
			continue
		}

		if lead.Type == card.Type {
			if lead.Number < card.Number {
				lead = card
			}
		} else {

			if card.isWhale() || card.isKraken() {
				lead = card
			}

			if (card.isParrot() || card.isMap() || card.isChest()) &&
				lead.isEscape() {
				lead = card
			}

			if (card.isRoger() || card.isMermaid() || card.isPirate() || card.isKing()) &&
				(lead.isStandardSuit() || lead.isEscape()) {
				lead = card
			}

			if card.isPirate() &&
				(lead.isRoger() || lead.isMermaid()) {
				lead = card
			}

			if card.isKing() &&
				(lead.isPirate() || lead.isRoger()) {
				lead = card
			}

			if card.isMermaid() &&
				lead.isKing() {
				lead = card
			}
		}
	}

	if lead.isWhale() {
		return suitLead.Id
	}

	if lead.isKraken() {
		return 0
	}

	if mermaidLead.Id != 0 && hasPirate && hasKing {
		return mermaidLead.Id
	}

	return lead.Id
}

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

// Set Owned user's cards
type Set struct {
	cards []Card
}

func (s Set) pickables(t Table) []CardId {
	var suit Card

	var options []CardId

	suit = t.suit()

	for _, card := range s.cards {
		if suit.Id == 0 || card.Type == suit.Type || card.isSpecial() {
			options = append(options, card.Id)
		}
	}

	return options
}
