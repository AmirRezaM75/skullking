package ws

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Initialize the game for given player
func initialize(game *Game, hub *Hub, receiverId int) {
	round := game.rounds[game.round]
	trick := round.tricks[game.trick]

	type Player struct {
		Id           int      `json:"id"`
		Username     string   `json:"username"`
		Avatar       string   `json:"avatar"`
		Score        int      `json:"score"`
		Bids         int      `json:"bids"`
		PickedCardId CardId   `json:"pickedCardId"`
		DealtCards   []CardId `json:"dealtCards"`
	}

	var players []Player

	for playerId, player := range game.players {
		var p Player

		p.Id = playerId
		p.Username = fmt.Sprintf("Username #%d", playerId)
		p.Avatar = player.avatar
		p.Score = round.scores[playerId]
		p.Bids = round.bids[playerId]
		p.PickedCardId = trick.pickedCards[playerId]

		// Receiver must not be aware of other cards
		if playerId == receiverId {
			p.DealtCards = round.dealtCards[playerId]
		}

		players = append(players, p)
	}

	content, _ := json.Marshal(struct {
		Round          int      `json:"round"`
		Trick          int      `json:"trick"`
		State          string   `json:"state"`
		ExpirationTime int      `json:"expirationTime"`
		PickingUserId  int      `json:"pickingUserId"`
		Players        []Player `json:"players"`
	}{
		Round:          game.round,
		Trick:          game.trick,
		State:          game.state,
		ExpirationTime: game.expirationTime,
		PickingUserId:  trick.pickingUserId,
		Players:        players,
	})

	m := &ServerMessage{
		ContentType: "json",
		Content:     string(content),
		Command:     CommandInit,
		gameId:      game.id,
		receiverId:  receiverId,
	}

	hub.dispatch <- m
}

func start(game *Game, hub *Hub) {
	var round = game.rounds[game.round]

	game.round++

	var deck Deck

	for _, card := range cards {
		deck.cards = append(deck.cards, card)
	}

	deck.shuffle()

	dealtCardIds := deck.deal(len(game.players), game.round)

	index := 0

	for _, player := range game.players {
		round.dealtCards[player.id] = dealtCardIds[index]

		playerCardIds, _ := json.Marshal(dealtCardIds[index])

		index++

		m := &ServerMessage{
			ContentType: "json",
			Content:     string(playerCardIds),
			Command:     CommandDeal,
			gameId:      game.id,
			receiverId:  player.id,
		}

		hub.dispatch <- m
	}

	content, _ := json.Marshal(
		BiddingStartedContent{
			Round:  game.round,
			EndsAt: time.Now().Add(WaitTime).Unix(),
		},
	)

	m := &ServerMessage{
		ContentType: "json",
		Content:     string(content),
		Command:     CommandBiddingStarted,
		gameId:      game.id,
	}

	hub.dispatch <- m

	game.state = StateBidding

	// TODO: Waiter
}

func pick(game *Game, hub *Hub, cardId int, senderId int) {
	var round = game.rounds[game.round]
	var trick = round.tricks[game.trick]

	if game.state != StatePicking || trick.pickingUserId == senderId {
		// Exception
		return
	}

	// TODO: Check if cardId is valid and exists in the last dealt cards

	trick.pickedCards[senderId] = CardId(cardId)

	var m = &ServerMessage{
		ContentType: "text",
		Content:     strconv.Itoa(cardId),
		Command:     CommandPicked,
		SenderId:    senderId,
		gameId:      game.id,
	}

	hub.dispatch <- m

	chooseNextPlayerForPicking(game, hub)
}

func endPicking(game *Game, hub *Hub) {
	pickForIdlePlayer(game, hub)

	chooseNextPlayerForPicking(game, hub)
}

func pickForIdlePlayer(game *Game, hub *Hub) {
	var round = game.rounds[game.round]
	var trick = round.tricks[game.trick]
	var pickerId = trick.pickingUserId

	if _, ok := trick.pickedCards[pickerId]; ok {
		return
	}

	trick.pickedCards[pickerId] = round.remainingCards[pickerId][0]

	content, _ := json.Marshal(struct {
		PlayerId int `json:"playerId"`
		CardId   int `json:"cardId"`
	}{
		PlayerId: pickerId,
		CardId:   int(trick.pickedCards[pickerId]),
	})

	m := &ServerMessage{
		ContentType: "json",
		Content:     string(content),
		Command:     CommandPicked,
		gameId:      game.id,
	}

	hub.dispatch <- m
}

func chooseNextPlayerForPicking(game *Game, hub *Hub) {
	var round = game.rounds[game.round]
	var trick = round.tricks[game.trick]

	content, _ := json.Marshal(
		PickingStartedContent{
			UserId: getNextPlayerIdForPicking(*game, trick),
			EndsAt: time.Now().Add(WaitTime).Unix(),
		},
	)

	m := &ServerMessage{
		ContentType: "json",
		Content:     string(content),
		Command:     CommandPickingStarted,
		gameId:      game.id,
	}

	hub.dispatch <- m

	// TODO: Timer
}
