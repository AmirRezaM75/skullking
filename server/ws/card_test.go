package ws

import (
	"reflect"
	"testing"
)

func TestHigherNumberOfSameTypeWins(t *testing.T) {
	CardIds := []CardId{Parrot1, Parrot4, Parrot2}
	cardId := winner(CardIds)
	if Parrot4 != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestHigherNumberOfLeadCardWins(t *testing.T) {
	CardIds := []CardId{Parrot1, Parrot4, Map12}
	cardId := winner(CardIds)
	if Parrot4 != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestWinnerCardId3(t *testing.T) {
	CardIds := []CardId{Parrot3, Map3, Parrot4, Chest10}
	cardId := winner(CardIds)
	if Parrot4 != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestWinnerCardId4(t *testing.T) {
	CardIds := []CardId{Map14, Map3, Roger1, Chest10}
	cardId := winner(CardIds)
	if Roger1 != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestPirateWinsSuitCards(t *testing.T) {
	CardIds := []CardId{Map13, Chest2, Roger1, Pirate2, Parrot6}
	cardId := winner(CardIds)
	if Pirate2 != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestPirateWinsMermaid(t *testing.T) {
	CardIds := []CardId{Map13, Mermaid2, Roger1, Pirate2, Parrot6}
	cardId := winner(CardIds)
	if Pirate2 != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestMermaidWinsSuitCard(t *testing.T) {
	CardIds := []CardId{Map13, Mermaid2, Roger1, Parrot6}
	cardId := winner(CardIds)
	if Mermaid2 != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestSkullKingWinsPirate(t *testing.T) {
	CardIds := []CardId{Map13, Pirate2, Roger1, SkullKing, Pirate4}
	cardId := winner(CardIds)
	if SkullKing != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestMermaidWinsSkullKing(t *testing.T) {
	CardIds := []CardId{Map13, Pirate2, Roger1, SkullKing, Pirate4, Mermaid1}
	cardId := winner(CardIds)
	if Mermaid1 != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestFirstEscapeCardWinsWhenAllCardsAreEscape(t *testing.T) {
	CardIds := []CardId{Escape3, Escape1, Escape2, Escape5}
	cardId := winner(CardIds)
	if Escape3 != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestEscapeCardsNeverWin(t *testing.T) {
	CardIds := []CardId{Chest1, Escape1, Escape2, Escape5}
	cardId := winner(CardIds)
	if Chest1 != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestEscapeCardsNeverWin2(t *testing.T) {
	CardIds := []CardId{Escape1, Chest1, Escape2, Escape5}
	cardId := winner(CardIds)
	if Chest1 != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestEscapeCardsNeverWin3(t *testing.T) {
	CardIds := []CardId{Escape1, Chest1, Escape2, Pirate2}
	cardId := winner(CardIds)
	if Pirate2 != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestEscapeCardsNeverWin4(t *testing.T) {
	CardIds := []CardId{Escape1, Chest1, Escape2, Chest2, Escape3}
	cardId := winner(CardIds)
	if Chest2 != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestHigherNumberWinsWhenWhaleIsLead(t *testing.T) {
	CardIds := []CardId{Escape1, Chest2, Escape2, Whale, Roger1, Pirate2}
	cardId := winner(CardIds)
	if Chest2 != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestNoWinnerWhenWhaleCardIsLeadAndThereIsNoSuitCard(t *testing.T) {
	CardIds := []CardId{Pirate1, Pirate2, Escape2, Mermaid2, SkullKing, Pirate3, Whale}
	cardId := winner(CardIds)
	if 0 != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestNoWinnerWhenKrakenCardIsLead(t *testing.T) {
	CardIds := []CardId{Pirate1, Mermaid1, Escape1, SkullKing, Kraken, Parrot1}
	cardId := winner(CardIds)
	if 0 != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestWhaleCardIsLeadWhenKrakenIsPickedBefore(t *testing.T) {
	CardIds := []CardId{Pirate1, Mermaid1, Escape1, SkullKing, Kraken, Parrot1, Whale}
	cardId := winner(CardIds)
	if Parrot1 != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestMermaidWinsIfPirateAndSkullKingAreAllPlayed1(t *testing.T) {
	CardIds := []CardId{Pirate1, Mermaid1, Escape1, SkullKing, Parrot1}
	cardId := winner(CardIds)
	if Mermaid1 != cardId {
		t.Errorf("Wrong winner card (%d).", cardId)
	}
}

func TestMermaidWinsIfPirateAndSkullKingAreAllPlayed2(t *testing.T) {
	CardIds := []CardId{Mermaid1, Pirate1, Escape1, SkullKing, Parrot1}
	cardId := winner(CardIds)
	if Mermaid1 != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestSkullKingWinsJollyRoger(t *testing.T) {
	CardIds := []CardId{Parrot3, Roger2, SkullKing, Parrot1}
	cardId := winner(CardIds)
	if SkullKing != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestPickableCardWhenTableIsEmpty(t *testing.T) {
	var table Table
	var set Set
	var card Card

	set.cards = []Card{
		card.fromId(Parrot2),
		card.fromId(Roger5),
		card.fromId(SkullKing),
		card.fromId(Map3),
		card.fromId(Map2),
	}

	expected := []CardId{
		Parrot2,
		Roger5,
		SkullKing,
		Map3,
		Map2,
	}

	options := set.pickables(table)

	if !reflect.DeepEqual(expected, options) {
		t.Error("Wrong pickable cards.")
	}
}

func TestPickableCardWhenFirstCardOnTableIsSuit(t *testing.T) {
	var table Table
	var set Set
	var card Card

	table.cards = []Card{
		card.fromId(Parrot3),
	}

	set.cards = []Card{
		card.fromId(Parrot2),
		card.fromId(Roger5),
		card.fromId(SkullKing),
		card.fromId(Map3),
		card.fromId(Map2),
	}

	expected := []CardId{
		Parrot2,
		SkullKing,
	}

	options := set.pickables(table)

	if !reflect.DeepEqual(expected, options) {
		t.Error("Wrong pickable cards.")
	}
}

func TestUserCanPickAnyCardIfNoCardMatchesTheSuit(t *testing.T) {
	var table Table
	var set Set
	var card Card

	table.cards = []Card{
		card.fromId(Parrot3),
	}

	set.cards = []Card{
		card.fromId(Chest1),
		card.fromId(Roger5),
		card.fromId(Pirate1),
		card.fromId(Map3),
		card.fromId(Map2),
	}

	expected := []CardId{
		Chest1,
		Roger5,
		Pirate1,
		Map3,
		Map2,
	}

	options := set.pickables(table)

	if !reflect.DeepEqual(expected, options) {
		t.Error("Wrong pickable cards.")
	}
}
