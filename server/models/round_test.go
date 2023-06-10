package models

import (
	"testing"
)

func TestCalculateScore(t *testing.T) {
	var trick = &Trick{
		Number:             1,
		WinnerPlayerId:     "Bilbo",
		StarterPlayerIndex: 0,
	}

	var round = Round{
		Bids: map[string]int{
			"Bilbo": 1,
			"Arwen": 0,
		},
		Number: 1,
		Tricks: []*Trick{trick},
		Scores: make(map[string]int, 2),
	}

	round.calculateScores()

	if round.Scores["Bilbo"] != 20 {
		t.Errorf("Expected 20 for Bilbo score, got %d", round.Scores["Bilbo"])
	}

	if round.Scores["Arwen"] != 10 {
		t.Errorf("Expected 10 for Arwen score, got %d", round.Scores["Arwen"])
	}
}

func TestCalculateScore2(t *testing.T) {

	var trick1 = &Trick{
		Number:             1,
		WinnerPlayerId:     "Bilbo",
		StarterPlayerIndex: 0,
	}
	var trick2 = &Trick{
		Number:             1,
		WinnerPlayerId:     "Bilbo",
		StarterPlayerIndex: 0,
	}

	var round = Round{
		Bids: map[string]int{
			"Bilbo": 2,
			"Arwen": 0,
		},
		Number: 2,
		Tricks: []*Trick{trick1, trick2},
		Scores: make(map[string]int, 2),
	}

	round.calculateScores()

	if round.Scores["Bilbo"] != 40 {
		t.Errorf("Expected 40 for Bilbo score, got %d", round.Scores["Bilbo"])
	}

	if round.Scores["Arwen"] != 20 {
		t.Errorf("Expected 20 for Arwen score, got %d", round.Scores["Arwen"])
	}
}

func TestCalculateScore3(t *testing.T) {
	var round = Round{
		Bids: map[string]int{
			"Bilbo": 0,
			"Arwen": 1,
		},
		Number: 2,
		Tricks: []*Trick{
			{
				WinnerPlayerId: "Bilbo",
			}, {
				WinnerPlayerId: "Bilbo",
			},
		},
		Scores: make(map[string]int, 2),
	}

	round.calculateScores()

	if round.Scores["Bilbo"] != -20 {
		t.Errorf("Expected -20 for Bilbo score, got %d", round.Scores["Bilbo"])
	}

	if round.Scores["Arwen"] != -10 {
		t.Errorf("Expected -10 for Arwen score, got %d", round.Scores["Arwen"])
	}
}

func TestCalculateScore4(t *testing.T) {
	var round = Round{
		Bids: map[string]int{
			"Bilbo": 0,
		},
		Number: 9,
		Tricks: []*Trick{
			{
				WinnerPlayerId: "Bilbo",
			}, {
				WinnerPlayerId: "Bilbo",
			},
		},
		Scores: make(map[string]int, 1),
	}

	round.calculateScores()

	if round.Scores["Bilbo"] != -90 {
		t.Errorf("Expected -90 for Bilbo score, got %d", round.Scores["Bilbo"])
	}
}

func TestCalculateScore5(t *testing.T) {
	var round = Round{
		Bids: map[string]int{
			"Bilbo": 0,
		},
		Number: 7,
		Tricks: []*Trick{},
		Scores: make(map[string]int, 1),
	}

	round.calculateScores()

	if round.Scores["Bilbo"] != 70 {
		t.Errorf("Expected 70 for Bilbo score, got %d", round.Scores["Bilbo"])
	}
}

func TestCalculateScoreWith14CardsBonusPoint(t *testing.T) {
	var round = Round{
		Bids: map[string]int{
			"Frodo":   3,
			"Gandalf": 1,
		},
		Number: 3,
		Tricks: []*Trick{
			{
				WinnerPlayerId: "Frodo",
				PickedCards: []PickedCard{
					{CardId: Parrot14},
					{CardId: Parrot1},
				},
			},
			{
				WinnerPlayerId: "Frodo",
				PickedCards: []PickedCard{
					{CardId: Roger14},
					{CardId: Escape1},
				},
			},
			{
				WinnerPlayerId: "Frodo",
				PickedCards: []PickedCard{
					{CardId: Parrot1},
					{CardId: Chest1},
				},
			},
		},
		Scores: make(map[string]int, 1),
	}

	round.calculateScores()

	if round.Scores["Frodo"] != 3*20+10+20 {
		t.Errorf("Expected 90 for Bilbo score, got %d", round.Scores["Frodo"])
	}
}

func TestCalculateScoreWithSkullKingBonusPoint(t *testing.T) {
	var round = Round{
		Bids: map[string]int{
			"Frodo":   3,
			"Gandalf": 1,
		},
		Number: 3,
		Tricks: []*Trick{
			{
				WinnerPlayerId: "Frodo",
				PickedCards: []PickedCard{
					{CardId: Parrot2},
					{CardId: Parrot1},
				},
			},
			{
				WinnerPlayerId: "Frodo",
				PickedCards: []PickedCard{
					{CardId: Chest1},
					{CardId: Parrot2},
				},
			},
			{
				WinnerCardId:   SkullKing,
				WinnerPlayerId: "Frodo",
				PickedCards: []PickedCard{
					{CardId: SkullKing},
					{CardId: Pirate1},
				},
			},
		},
		Scores: make(map[string]int, 1),
	}

	round.calculateScores()

	if round.Scores["Frodo"] != 60+30 {
		t.Errorf("Expected 90 for Bilbo score, got %d", round.Scores["Frodo"])
	}
}

func TestNoBonusPointWhenYouBidWrong(t *testing.T) {
	var round = Round{
		Bids: map[string]int{
			"Frodo":   3,
			"Gandalf": 1,
		},
		Number: 3,
		Tricks: []*Trick{
			{
				WinnerPlayerId: "Frodo",
				PickedCards: []PickedCard{
					{CardId: Parrot2},
					{CardId: Parrot1},
				},
			},
			{
				WinnerPlayerId: "",
				PickedCards: []PickedCard{
					{CardId: Kraken},
					{CardId: Parrot2},
				},
			},
			{
				WinnerPlayerId: "Frodo",
				PickedCards: []PickedCard{
					{CardId: SkullKing},
					{CardId: Pirate1},
				},
			},
		},
		Scores: make(map[string]int, 1),
	}

	round.calculateScores()

	if round.Scores["Frodo"] != -10 {
		t.Errorf("Expected -10 for Bilbo score, got %d", round.Scores["Frodo"])
	}
}
