package models

import "testing"

func TestGetWinnerBonusPointForEachFourteenCards1(t *testing.T) {
	var trick = Trick{
		PickedCards: []PickedCard{
			{CardId: Parrot2},
			{CardId: Chest14},
		},
	}

	bonus := trick.getWinnerBonusPoint()

	if bonus != 10 {
		t.Errorf("Expected 10 for bonus point, got %d", bonus)
	}
}

func TestGetWinnerBonusPointForEachFourteenCards2(t *testing.T) {
	var trick = Trick{
		PickedCards: []PickedCard{
			{CardId: Parrot14},
			{CardId: Chest14},
			{CardId: Map14},
			{CardId: Parrot2},
		},
	}

	bonus := trick.getWinnerBonusPoint()

	if bonus != 30 {
		t.Errorf("Expected 30 for bonus point, got %d", bonus)
	}
}

func TestGetWinnerBonusPointForEachFourteenCards3(t *testing.T) {
	var trick = Trick{
		PickedCards: []PickedCard{
			{CardId: Parrot14},
			{CardId: Chest4},
			{CardId: Roger14},
			{CardId: Parrot2},
		},
	}

	bonus := trick.getWinnerBonusPoint()

	if bonus != 20+10 {
		t.Errorf("Expected 30 for bonus point, got %d", bonus)
	}
}

func TestGetWinnerBonusPointForEachMermaidTakenByAPirate(t *testing.T) {
	var trick = Trick{
		PickedCards: []PickedCard{
			{CardId: Mermaid1},
			{CardId: Pirate2},
			{CardId: Mermaid2},
			{CardId: Parrot2},
		},
		WinnerCardId: Pirate2,
	}

	bonus := trick.getWinnerBonusPoint()

	if bonus != 20+20 {
		t.Errorf("Expected 40 for bonus point, got %d", bonus)
	}
}

func TestGetWinnerBonusPointForEachPirateTakenBySkullKing(t *testing.T) {
	var trick = Trick{
		PickedCards: []PickedCard{
			{CardId: Pirate1},
			{CardId: Pirate2},
			{CardId: SkullKing},
			{CardId: Pirate3},
		},
		WinnerCardId: SkullKing,
	}

	bonus := trick.getWinnerBonusPoint()

	if bonus != 3*30 {
		t.Errorf("Expected 90 for bonus point, got %d", bonus)
	}
}

func TestGetWinnerBonusPointForTakingTheSkullKingWithAMermaid(t *testing.T) {
	var trick = Trick{
		PickedCards: []PickedCard{
			{CardId: Pirate1},
			{CardId: SkullKing},
			{CardId: Mermaid1},
		},
		WinnerCardId: Mermaid1,
	}

	bonus := trick.getWinnerBonusPoint()

	if bonus != 40 {
		t.Errorf("Expected 40 for bonus point, got %d", bonus)
	}
}

func TestGetWinnerBonusPointForTakingTheSkullKingWithAMermaidIncludingFourteenCard(t *testing.T) {
	var trick = Trick{
		PickedCards: []PickedCard{
			{CardId: Chest14},
			{CardId: Pirate1},
			{CardId: SkullKing},
			{CardId: Mermaid1},
		},
		WinnerCardId: Mermaid1,
	}

	bonus := trick.getWinnerBonusPoint()

	if bonus != 10+40 {
		t.Errorf("Expected 50 for bonus point, got %d", bonus)
	}
}
