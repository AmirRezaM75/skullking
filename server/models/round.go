package models

import (
	"github.com/AmirRezaM75/skull-king/pkg/support"
)

type Round struct {
	Number             int
	DealtCards         map[string][]CardId
	Bids               map[string]int
	Tricks             []*Trick
	StarterPlayerIndex int
	Scores             map[string]int
}

func (round Round) getDealtCardIdsByPlayerId(playerId string) []uint16 {
	var cardIds []uint16

	for _, cardId := range round.DealtCards[playerId] {
		cardIds = append(cardIds, uint16(cardId))
	}

	return cardIds
}

func (round Round) getPickedCardIdsByPlayerId(playerId string) []CardId {
	var cardIds []CardId

	for _, trick := range round.Tricks {
		if trick == nil {
			break
		}
		for _, pickedCard := range trick.PickedCards {
			if playerId == pickedCard.PlayerId {
				cardIds = append(cardIds, pickedCard.CardId)
			}
		}
	}

	return cardIds
}

func (round *Round) calculateScores() {
	var wonTricks = make(map[string]int, len(round.Bids))

	for playerId, _ := range round.Bids {
		if _, ok := wonTricks[playerId]; !ok {
			wonTricks[playerId] = 0
		}

		for _, trick := range round.Tricks {
			if trick.WinnerPlayerId == playerId {
				wonTricks[trick.WinnerPlayerId] += 1
			}
		}
	}

	for playerId, bid := range round.Bids {
		if wonTricks[playerId] != bid && bid != 0 {
			diff := support.Abs(wonTricks[playerId] - bid)
			round.Scores[playerId] = -10 * diff
		}

		if wonTricks[playerId] != bid && bid == 0 {
			round.Scores[playerId] = -10 * round.Number
		}

		if wonTricks[playerId] == bid && bid != 0 {
			round.Scores[playerId] = 20*bid + round.getBonusPointByPlayerId(playerId)
		}

		if wonTricks[playerId] == bid && bid == 0 {
			round.Scores[playerId] = 10 * round.Number
		}
	}
}

func (round Round) getBonusPointByPlayerId(playerId string) int {
	var bonus = 0

	for _, trick := range round.Tricks {
		if trick.WinnerPlayerId == playerId {
			bonus += trick.getWinnerBonusPoint()
		}
	}

	return bonus
}

func (round Round) getRemainingCardIds(playerId string) []CardId {
	var remainingCardIds []CardId

	pickedCardIds := round.getPickedCardIdsByPlayerId(playerId)

outerLoop:
	for _, dealtCardId := range round.DealtCards[playerId] {
		for _, pickedCardId := range pickedCardIds {
			if pickedCardId == dealtCardId {
				continue outerLoop
			}
		}
		remainingCardIds = append(remainingCardIds, dealtCardId)
	}

	return remainingCardIds
}

func (round Round) getRemainingIntCardIds(playerId string) []uint16 {
	var remainingCardIds []uint16

	pickedCardIds := round.getPickedCardIdsByPlayerId(playerId)

outerLoop:
	for _, dealtCardId := range round.DealtCards[playerId] {
		for _, pickedCardId := range pickedCardIds {
			if pickedCardId == dealtCardId {
				continue outerLoop
			}
		}
		remainingCardIds = append(remainingCardIds, uint16(dealtCardId))
	}

	return remainingCardIds
}
