package models

import (
	"encoding/json"
	"fmt"
	"github.com/AmirRezaM75/skull-king/constants"
	"github.com/AmirRezaM75/skull-king/pkg/support"
	"strconv"
	"time"
)

type Game struct {
	Id             string
	Round          int
	Trick          int
	State          string
	ExpirationTime int
	Players        map[string]*Player
	Rounds         map[int]*Round
}

func (game *Game) GetNextPlayerIdForPicking(trick Trick) string {
	var currentPickingPlayerHaveFound = false

	var pickerId string

	for playerId, _ := range game.Players {
		if currentPickingPlayerHaveFound {
			pickerId = playerId
		}

		if playerId == trick.PickingUserId {
			currentPickingPlayerHaveFound = true
		}
	}

	return pickerId
}

func (game *Game) Start(hub *Hub) {
	var round = game.Rounds[game.Round]

	game.Round++

	var deck Deck

	for _, card := range Cards {
		deck.Cards = append(deck.Cards, card)
	}

	deck.Shuffle()

	dealtCardIds := deck.Deal(len(game.Players), game.Round)

	index := 0

	for _, player := range game.Players {
		round.DealtCards[player.Id] = dealtCardIds[index]

		playerCardIds, _ := json.Marshal(dealtCardIds[index])

		index++

		m := &ServerMessage{
			Content:    string(playerCardIds),
			Command:    constants.CommandDeal,
			GameId:     game.Id,
			ReceiverId: player.Id,
		}

		hub.Dispatch <- m
	}

	content, _ := json.Marshal(
		struct {
			Round  int   `json:"round"`
			EndsAt int64 `json:"endsAt"`
		}{
			Round:  game.Round,
			EndsAt: time.Now().Add(constants.WaitTime).Unix(),
		},
	)

	m := &ServerMessage{
		Content: string(content),
		Command: constants.CommandBiddingStarted,
		GameId:  game.Id,
	}

	hub.Dispatch <- m

	game.State = constants.StateBidding

	// TODO: Waiter
}

func (game *Game) EndPicking(hub *Hub) {
	pickForIdlePlayer(game, hub)

	chooseNextPlayerForPicking(game, hub)
}

func pickForIdlePlayer(game *Game, hub *Hub) {
	var round = game.Rounds[game.Round]
	var trick = round.Tricks[game.Trick]
	var pickerId = trick.PickingUserId

	if _, ok := trick.PickedCards[pickerId]; ok {
		return
	}

	trick.PickedCards[pickerId] = round.RemainingCards[pickerId][0]

	content, _ := json.Marshal(struct {
		PlayerId string `json:"playerId"`
		CardId   int    `json:"cardId"`
	}{
		PlayerId: pickerId,
		CardId:   int(trick.PickedCards[pickerId]),
	})

	m := &ServerMessage{
		Content: string(content),
		Command: constants.CommandPicked,
		GameId:  game.Id,
	}

	hub.Dispatch <- m
}

func chooseNextPlayerForPicking(game *Game, hub *Hub) {
	var round = game.Rounds[game.Round]
	var trick = round.Tricks[game.Trick]

	content, _ := json.Marshal(
		struct {
			UserId string `json:"userId"`
			EndsAt int64  `json:"endsAt"`
		}{
			UserId: game.GetNextPlayerIdForPicking(trick),
			EndsAt: time.Now().Add(constants.WaitTime).Unix(),
		},
	)

	m := &ServerMessage{
		Content: string(content),
		Command: constants.CommandPickingStarted,
		GameId:  game.Id,
	}

	hub.Dispatch <- m

	// TODO: Timer
}

func (game *Game) Initialize(hub *Hub, receiverId string) {

	round := game.Rounds[game.Round]

	type Player struct {
		Id           string   `json:"id"`
		Username     string   `json:"username"`
		Avatar       string   `json:"avatar"`
		Score        int      `json:"score"`
		Bids         int      `json:"bids"`
		PickedCardId CardId   `json:"pickedCardId"`
		DealtCards   []CardId `json:"dealtCards"`
	}

	var players []Player

	for playerId, player := range game.Players {
		var p Player

		p.Id = playerId
		p.Username = player.Username
		p.Avatar = player.Avatar

		if round != nil {
			p.Score = round.Scores[playerId]
			p.Bids = round.Bids[playerId]
			p.PickedCardId = round.Tricks[game.Trick].PickedCards[playerId]

			// Receiver must not be aware of other cards
			if playerId == receiverId {
				p.DealtCards = round.DealtCards[playerId]
			}
		}

		players = append(players, p)
	}

	content := struct {
		Round          int      `json:"round"`
		Trick          int      `json:"trick"`
		State          string   `json:"state"`
		ExpirationTime int      `json:"expirationTime"`
		PickingUserId  string   `json:"pickingUserId"`
		Players        []Player `json:"players"`
	}{
		Round:          game.Round,
		Trick:          game.Trick,
		State:          game.State,
		ExpirationTime: game.ExpirationTime,
		Players:        players,
	}

	if round != nil {
		content.PickingUserId = round.Tricks[game.Trick].PickingUserId
	}

	m := &ServerMessage{
		Content:    content,
		Command:    constants.CommandInit,
		GameId:     game.Id,
		ReceiverId: receiverId,
	}

	hub.Dispatch <- m
}

func (game *Game) Pick(hub *Hub, cardId int, senderId string) {
	var round = game.Rounds[game.Round]
	var trick = round.Tricks[game.Trick]

	if game.State != constants.StatePicking || trick.PickingUserId == senderId {
		// TODO: Exception
		return
	}

	// TODO: Check if cardId is valid and exists in the last dealt cards

	trick.PickedCards[senderId] = CardId(cardId)

	var m = &ServerMessage{
		Content:  strconv.Itoa(cardId),
		Command:  constants.CommandPicked,
		SenderId: senderId,
		GameId:   game.Id,
	}

	hub.Dispatch <- m

	chooseNextPlayerForPicking(game, hub)
}

func (game *Game) GetAvailableAvatar() string {
outerLoop:
	for _, number := range support.Fill(constants.MaxPlayers) {
		for _, player := range game.Players {
			if player.Avatar == fmt.Sprintf("avatar-%d", number) {
				continue outerLoop
			}
		}
		return fmt.Sprintf("avatar-%d", number)
	}

	return ""
}

func (game *Game) Join(hub *Hub, player *Player) {
	// Must send JOINED command after INIT command
	// Because it preserves order of players in frontend
	// TODO: Or we can stop sending JOINED to the new joiner player
	m := &ServerMessage{
		Command: constants.CommandJoined,
		Content: struct {
			Id       string `json:"id"`
			Username string `json:"username"`
			Avatar   string `json:"avatar"`
		}{
			Id:       player.Id,
			Username: player.Username,
			Avatar:   player.Avatar,
		},
		GameId:   player.GameId,
		SenderId: player.Id,
	}

	hub.Dispatch <- m
}

func (game *Game) Left(hub *Hub, playerId string) {
	m := &ServerMessage{
		Content: struct {
			PlayerId string `json:"playerId"`
		}{PlayerId: playerId},
		Command:  constants.CommandLeft,
		GameId:   game.Id,
		SenderId: playerId,
	}

	hub.Dispatch <- m
}
