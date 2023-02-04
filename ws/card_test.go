package ws

import "testing"

func TestWinnerCardId1(t *testing.T) {
	CardIds := []CardId{Parrot1, Parrot4, Parrot2}
	cardId := winner(CardIds)
	if Parrot4 != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestWinnerCardId2(t *testing.T) {
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

func TestWinnerCardId5(t *testing.T) {
	CardIds := []CardId{Map13, Chest2, Roger1, Pirate2, Parrot6}
	cardId := winner(CardIds)
	if Pirate2 != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestWinnerCardId6(t *testing.T) {
	CardIds := []CardId{Map13, Mermaid2, Roger1, Pirate2, Parrot6}
	cardId := winner(CardIds)
	if Pirate2 != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestWinnerCardId7(t *testing.T) {
	CardIds := []CardId{Map13, Mermaid2, Roger1, Parrot6}
	cardId := winner(CardIds)
	if Mermaid2 != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestWinnerCardId8(t *testing.T) {
	CardIds := []CardId{Map13, Pirate2, Roger1, SkullKing, Pirate4}
	cardId := winner(CardIds)
	if SkullKing != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestWinnerCardId9(t *testing.T) {
	CardIds := []CardId{Map13, Pirate2, Roger1, SkullKing, Pirate4, Mermaid1}
	cardId := winner(CardIds)
	if Mermaid1 != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestWinnerCardId10(t *testing.T) {
	CardIds := []CardId{Escape3, Escape1, Escape2, Escape5}
	cardId := winner(CardIds)
	if Escape3 != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestWinnerCardId11(t *testing.T) {
	CardIds := []CardId{Chest1, Escape1, Escape2, Escape5}
	cardId := winner(CardIds)
	if Chest1 != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestWinnerCardId12(t *testing.T) {
	CardIds := []CardId{Escape1, Chest1, Escape2, Escape5}
	cardId := winner(CardIds)
	if Chest1 != cardId {
		t.Errorf("Wrong winner card.")
	}
}

func TestWinnerCardId13(t *testing.T) {
	CardIds := []CardId{Escape1, Chest1, Escape2, Pirate2}
	cardId := winner(CardIds)
	if Pirate2 != cardId {
		t.Errorf("Wrong winner card.")
	}
}
