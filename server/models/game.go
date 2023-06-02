package models

import (
	"fmt"
	"github.com/AmirRezaM75/skull-king/constants"
	"github.com/AmirRezaM75/skull-king/pkg/support"
	"github.com/AmirRezaM75/skull-king/responses"
	"log"
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

func (game *Game) findPlayerIndexForPicking() int {
	pickedCardsCount := len(game.Rounds[game.Round].Tricks[game.Trick].PickedCards)

	if pickedCardsCount != 0 {
		index := game.Rounds[game.Round].Tricks[game.Trick].StarterPlayerIndex + pickedCardsCount
		if index > len(game.Players) {
			index = 1
		}
		return index
	}

	if game.Round == 1 && game.Trick == 1 {
		return 1
	}

	if game.Round > 1 && game.Trick == 1 {
		index := game.Rounds[game.Round-1].StarterPlayerIndex + 1
		if index > len(game.Players) {
			index = 1
		}
		return index
	}

	playerId := game.Rounds[game.Round].Tricks[game.Trick-1].WinnerPlayerId

	if playerId != "" {
		for _, player := range game.Players {
			if player.Id == playerId {
				return player.Index
			}
		}
	} else {
		return game.Rounds[game.Round].Tricks[game.Trick-1].StarterPlayerIndex
	}

	log.Fatalln("Unable to find playerId within for loop.")
	return 1
}

func (game *Game) setNextPlayerForPicking() string {
	index := game.findPlayerIndexForPicking()
	pickedCardsCount := len(game.Rounds[game.Round].Tricks[game.Trick].PickedCards)

	if pickedCardsCount == 0 {
		if game.Trick == 1 {
			game.Rounds[game.Round].StarterPlayerIndex = index
		}
		game.Rounds[game.Round].Tricks[game.Trick].StarterPlayerIndex = index
	}

	for _, player := range game.Players {
		if player.Index == index {
			game.Rounds[game.Round].Tricks[game.Trick].PickingUserId = player.Id
			return player.Id
		}
	}

	log.Fatalln("Unable to find playerId within for loop.")
	return ""
}

func (game *Game) NextRound(hub *Hub) {
	game.Round++
	game.Trick = 1
	game.State = constants.StateDealing

	var deck Deck

	deck.Shuffle()

	dealtCardIds := deck.Deal(len(game.Players), game.Round)

	round := Round{
		Number:     game.Round,
		Scores:     make(map[string]int, len(game.Players)),
		DealtCards: make(map[string][]CardId, len(game.Players)),
		Bids:       make(map[string]int, len(game.Players)),
		Tricks:     make(map[int]*Trick, game.Round),
	}

	index := 0

	for _, player := range game.Players {
		player.Index = index + 1

		round.DealtCards[player.Id] = dealtCardIds[index]
		round.Scores[player.Id] = 0
		round.Bids[player.Id] = 0

		index++
	}

	trick := &Trick{
		Number:      game.Trick,
		PickedCards: make(map[string]CardId, constants.MaxPlayers),
	}
	round.Tricks[game.Trick] = trick
	game.Rounds[game.Round] = &round

	for _, player := range game.Players {
		content := responses.DealResponse{
			Round: game.Round,
			Trick: game.Trick,
			Cards: round.getDealtCardIdsByPlayerId(player.Id),
			State: game.State,
		}

		m := &ServerMessage{
			Content:    content,
			Command:    constants.CommandDeal,
			GameId:     game.Id,
			ReceiverId: player.Id,
			SenderId:   "SERVER",
		}

		hub.Dispatch <- m
	}

	game.startBidding(hub)
}

func (game *Game) startBidding(hub *Hub) {
	duration := game.getBiddingExpirationDuration()

	m := &ServerMessage{
		Content: responses.StartBidding{
			EndsAt: time.Now().Add(duration).Unix(),
		},
		Command:  constants.CommandStartBidding,
		SenderId: "SERVER",
		GameId:   game.Id,
	}

	hub.Dispatch <- m

	game.State = constants.StateBidding

	timer := time.NewTimer(duration)

	go func() {
		<-timer.C
		game.endBidding(hub)
	}()
}

func (game *Game) endBidding(hub *Hub) {
	m := &ServerMessage{
		Content:  nil,
		Command:  constants.CommandEndBidding,
		SenderId: "SERVER",
		GameId:   game.Id,
	}

	hub.Dispatch <- m

	game.startPicking(hub)
}

func (game *Game) startPicking(hub *Hub) {
	playerId := game.setNextPlayerForPicking()

	if playerId == "" {
		log.Fatalln("No player id is found for picking")
		return
	}

	content := responses.StartPicking{
		PlayerId: playerId,
		EndsAt:   time.Now().Add(constants.WaitTime).Unix(),
	}

	m := &ServerMessage{
		Content: content,
		Command: constants.CommandStartPicking,
		GameId:  game.Id,
	}

	hub.Dispatch <- m

	timer := time.NewTimer(constants.WaitTime)
	go func() {
		<-timer.C
		game.endPicking(hub)
	}()
}

func (game *Game) endPicking(hub *Hub) {
	game.pickForIdlePlayer(hub)

	if game.isTrickOver() {
		game.announceTrickWinner(hub)
		game.nextTrick(hub)
	} else {
		game.startPicking(hub)
	}
}

func (game *Game) isTrickOver() bool {
	var round = game.Rounds[game.Round]
	var trick = round.Tricks[game.Trick]
	return len(trick.PickedCards) == len(game.Players)
}

func (game *Game) announceTrickWinner(hub *Hub) {
	cardId, playerId := game.findTrickWinner()

	game.Rounds[game.Round].Tricks[game.Trick].WinnerPlayerId = playerId

	content := responses.AnnounceTrickWinner{
		PlayerId: playerId,
		CardId:   int(cardId),
	}

	m := &ServerMessage{
		Content:  content,
		Command:  constants.CommandAnnounceTrickWinner,
		GameId:   game.Id,
		SenderId: "SERVER",
	}

	hub.Dispatch <- m
}

func (game *Game) nextTrick(hub *Hub) {
	if game.Trick == game.Round {
		game.NextRound(hub)
		return
	}

	game.Trick++
	game.Rounds[game.Round].Tricks[game.Trick] = &Trick{
		Number:      game.Trick,
		PickedCards: make(map[string]CardId, constants.MaxPlayers),
	}

	game.startPicking(hub)
}

func (game *Game) findTrickWinner() (CardId, string) {
	var round = game.Rounds[game.Round]
	var trick = round.Tricks[game.Trick]

	var cardIds []CardId
	for _, cardId := range trick.PickedCards {
		cardIds = append(cardIds, cardId)
	}

	winnerCardId := winner(cardIds)

	if winnerCardId == 0 {
		return winnerCardId, ""
	}

	var winnerPlayerId string
	for playerId, cardId := range trick.PickedCards {
		if cardId == winnerCardId {
			winnerPlayerId = playerId
			break
		}
	}

	return winnerCardId, winnerPlayerId
}

func (game *Game) pickForIdlePlayer(hub *Hub) {
	var round = game.Rounds[game.Round]
	var trick = round.Tricks[game.Trick]
	var pickerId = trick.PickingUserId

	if _, ok := trick.PickedCards[pickerId]; ok {
		return
	}

	remainingCards := game.getRemainingCardsForPlayerId(pickerId)
	trick.PickedCards[pickerId] = remainingCards[0]

	content := responses.Pick{
		PlayerId: pickerId,
		CardId:   int(trick.PickedCards[pickerId]),
	}

	m := &ServerMessage{
		Content: content,
		Command: constants.CommandPicked,
		GameId:  game.Id,
	}

	hub.Dispatch <- m
}

func (game *Game) getRemainingCardsForPlayerId(playerId string) []CardId {
	var remainingCardIds []CardId
	var round = game.Rounds[game.Round]
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

	game.startPicking(hub)
}

func (game *Game) GetAvailableAvatar() string {
	// TODO: Not working properly
outerLoop:
	for _, number := range support.Fill(constants.MaxPlayers) {
		for _, player := range game.Players {
			if player.Avatar == fmt.Sprintf("%d.jpg", number) {
				continue outerLoop
			}
		}
		return fmt.Sprintf("%d.jpg", number)
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

func (game *Game) getBiddingExpirationDuration() time.Duration {
	// As the round number increases, it takes more time to complete the card dealing animation.
	// Therefore, we need to increase the wait time for each level
	// Each animation takes about 2 seconds
	return constants.WaitTime + time.Duration(game.Round)*2*time.Second
}
