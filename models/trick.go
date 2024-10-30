package models

import "skullking/responses"

type Trick struct {
	Number int
	// This is useful to find out who is eligible when receiving 'PICK' command.
	PickingUserId string
	// We cannot use a map for picked cards because the order in which the cards are picked
	// affects the pickable cards. In Go, a map is not a suitable collection to maintain
	// the sequence of elements as it is unordered.
	PickedCards        []PickedCard
	WinnerPlayerId     string
	WinnerCardId       CardId
	StarterPlayerIndex int
}

type PickedCard struct {
	PlayerId string
	CardId   CardId
}

func (trick Trick) isPlayerPicked(playerId string) bool {
	for _, pickedCard := range trick.PickedCards {
		if playerId == pickedCard.PlayerId {
			return true
		}
	}

	return false
}

func (trick Trick) getAllPickedCards() []responses.TableCard {
	var tableCards []responses.TableCard

	for _, pickedCard := range trick.PickedCards {
		tableCards = append(tableCards, responses.TableCard{
			PlayerId: pickedCard.PlayerId,
			CardId:   uint16(pickedCard.CardId),
		})
	}

	return tableCards
}

func (trick Trick) getAllPickedCardIds() []CardId {
	var cardIds []CardId

	for _, pickedCard := range trick.PickedCards {
		cardIds = append(cardIds, pickedCard.CardId)
	}

	return cardIds
}

func (trick Trick) getAllPickedIntCardIds() []uint16 {
	// The SkullKing AI can't handle null
	var cardIds = make([]uint16, 0)

	for _, cardId := range trick.getAllPickedCardIds() {
		cardIds = append(cardIds, uint16(cardId))
	}

	return cardIds
}

func (trick Trick) getWinnerBonusPoint() int {
	var bonus int

	var winnerCard = newCardFromId(trick.WinnerCardId)

	for _, pickedCard := range trick.PickedCards {
		if pickedCard.CardId == Parrot14 ||
			pickedCard.CardId == Chest14 ||
			pickedCard.CardId == Map14 {
			bonus += 10
			continue
		}

		if pickedCard.CardId == Roger14 {
			bonus += 20
			continue
		}

		card := newCardFromId(pickedCard.CardId)

		if winnerCard.isPirate() && card.isMermaid() {
			bonus += 20
			continue
		}

		if winnerCard.isKing() && card.isPirate() {
			bonus += 30
			continue
		}

		if winnerCard.isMermaid() && card.isKing() {
			bonus += 40
		}
	}
	return bonus
}

func (trick Trick) getWinner() (CardId, string) {
	var cardIds []CardId

	for _, pickedCard := range trick.PickedCards {
		cardIds = append(cardIds, pickedCard.CardId)
	}

	winnerCardId := winner(cardIds)

	if winnerCardId == 0 {
		return winnerCardId, ""
	}

	var winnerPlayerId string
	for _, pickedCard := range trick.PickedCards {
		if pickedCard.CardId == winnerCardId {
			winnerPlayerId = pickedCard.PlayerId
			break
		}
	}

	return winnerCardId, winnerPlayerId
}
