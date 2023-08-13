package models

import (
	"errors"
	"fmt"
	"github.com/AmirRezaM75/skull-king/constants"
	"github.com/AmirRezaM75/skull-king/pkg/support"
	"github.com/AmirRezaM75/skull-king/responses"
	"log"
	"sort"
	"time"
)

type Game struct {
	Id             string
	Round          int
	Trick          int
	State          string
	ExpirationTime int64
	Players        map[string]*Player
	Rounds         [constants.MaxRounds]*Round
	CreatorId      string
	CreatedAt      int64
}

func (game *Game) findPlayerIndexForPicking() int {
	pickedCardsCount := len(game.getCurrentTrick().PickedCards)

	if pickedCardsCount != 0 {
		index := game.getCurrentTrick().StarterPlayerIndex + pickedCardsCount
		if index > len(game.Players) {
			index -= len(game.Players)
		}
		return index
	}

	if game.Round == 1 && game.Trick == 1 {
		return 1
	}

	if game.Round > 1 && game.Trick == 1 {
		index := game.getPreviousRound().StarterPlayerIndex + 1
		if index > len(game.Players) {
			index = 1
		}
		return index
	}

	playerId := game.getPreviousTrick().WinnerPlayerId

	if playerId != "" {
		for _, player := range game.Players {
			if player.Id == playerId {
				return player.Index
			}
		}
	} else {
		// TODO: Kraken - The next trick is led by the player who would have won the trick.
		// TODO: Whale - The person who played the White Whale is the next to lead.
		return game.getPreviousTrick().StarterPlayerIndex
	}

	log.Println(
		fmt.Sprintf(
			"Unable to find playerId within for loop. [gameId: %s, round: %d, trick: %d]",
			game.Id,
			game.Round,
			game.Trick,
		),
	)

	return 1
}

func (game *Game) setNextPlayerForPicking() string {
	index := game.findPlayerIndexForPicking()
	pickedCardsCount := len(game.getCurrentTrick().PickedCards)

	if pickedCardsCount == 0 {
		if game.Trick == 1 {
			game.getCurrentRound().StarterPlayerIndex = index
		}
		game.getCurrentTrick().StarterPlayerIndex = index
	}

	for _, player := range game.Players {
		if player.Index == index {
			game.getCurrentTrick().PickingUserId = player.Id
			return player.Id
		}
	}

	log.Println(fmt.Sprintf("Unable to find player by index %d", index))

	return ""
}

func (game *Game) endGame(hub *Hub) {
	m := &ServerMessage{
		Command: constants.CommandEndGame,
		GameId:  game.Id,
	}

	hub.Dispatch <- m

	_, err := hub.GameRepository.Create(*game)

	if err != nil {
		fmt.Println("Can't persist game in database", err.Error())
	}

	delete(hub.Games, game.Id)
}

func (game *Game) NextRound(hub *Hub) {
	if game.Round == constants.MaxRounds {
		game.endGame(hub)
		return
	}

	game.Round++
	game.Trick = 1
	game.State = constants.StateDealing

	var deck Deck

	deck.Shuffle()

	dealtCardIds := deck.Deal(len(game.Players), game.Round)

	round := Round{
		Number:     game.Round,
		DealtCards: make(map[string][]CardId, len(game.Players)),
		Bids:       make(map[string]int, len(game.Players)),
		Tricks:     make([]*Trick, game.Round),
		Scores:     make(map[string]int, len(game.Players)),
	}

	index := 0

	for _, player := range game.getPlayers() {
		// Initially, we assign a Unix time as an index for each player when they join.
		// However, we require a sequential index starting from 1 to identify the next player for picking.
		if player.Index > constants.MaxPlayers {
			player.Index = index + 1
		}

		round.DealtCards[player.Id] = dealtCardIds[index]
		round.Bids[player.Id] = 0

		index++
	}

	trick := &Trick{
		Number: game.Trick,
	}
	round.Tricks[game.Trick-1] = trick
	game.Rounds[game.Round-1] = &round

	for _, player := range game.Players {
		content := responses.Deal{
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
		}

		hub.Dispatch <- m
	}

	game.startBidding(hub)
}

func (game *Game) startBidding(hub *Hub) {
	duration := game.getBiddingExpirationDuration()

	game.State = constants.StateBidding

	content := responses.StartBidding{
		EndsAt: time.Now().Add(duration).Unix(),
		State:  game.State,
		Round:  game.Round,
	}

	m := &ServerMessage{
		Content: content,
		Command: constants.CommandStartBidding,
		GameId:  game.Id,
	}

	hub.Dispatch <- m

	game.ExpirationTime = content.EndsAt

	timer := time.NewTimer(duration)

	go func() {
		<-timer.C
		game.endBidding(hub)
	}()
}

func (game *Game) endBidding(hub *Hub) {
	game.ExpirationTime = 0
	game.State = "" // TODO: Better name!?

	var bids []responses.Bid

	for playerId, number := range game.getCurrentRound().Bids {
		bids = append(bids, responses.Bid{
			PlayerId: playerId,
			Number:   number,
		})
	}

	content := responses.EndBidding{Bids: bids}

	m := &ServerMessage{
		Content: content,
		Command: constants.CommandEndBidding,
		GameId:  game.Id,
	}

	hub.Dispatch <- m

	game.startPicking(hub)
}

func (game *Game) startPicking(hub *Hub) {
	game.State = constants.StatePicking

	pickerId := game.setNextPlayerForPicking()

	if pickerId == "" {
		log.Println("No player id is found for picking")
	}

	content := responses.StartPicking{
		PlayerId: pickerId,
		EndsAt:   game.getPickingExpirationTime(),
		CardIds:  []uint16{},
		State:    game.State,
	}

	for _, player := range game.Players {
		if pickerId == player.Id {
			content.CardIds = game.getPickableIntCardIds(pickerId)
		}

		m := &ServerMessage{
			Content:    content,
			Command:    constants.CommandStartPicking,
			GameId:     game.Id,
			ReceiverId: player.Id,
		}

		hub.Dispatch <- m
	}

	game.ExpirationTime = content.EndsAt

	var trick = game.getCurrentTrick()

	timer := time.NewTimer(
		time.Unix(content.EndsAt, 0).Sub(time.Now()),
	)
	go func() {
		<-timer.C
		game.stopPicking(hub, pickerId, trick)
	}()
}

func (game *Game) getPickingExpirationTime() int64 {
	t := constants.WaitTime

	var trick = game.getCurrentTrick()

	// We need 4 seconds extra time to make sure all animations are completed in client side
	// 2 seconds waiting for announcing trick winner + 2 seconds for picked card animation
	if trick.Number != 1 && len(trick.PickedCards) == 0 {
		t += time.Second * 4
	}

	return time.Now().Add(t).Unix()
}

// stopPicking needs to get trick as parameter because
// the trick might be increased when this function is called.
func (game *Game) stopPicking(hub *Hub, playerId string, trick *Trick) {
	// When picking time is expired there is no need to take any further action
	// if player already picked the card because we already called endPicking function
	if !trick.isPlayerPicked(playerId) {
		game.pickForIdlePlayer(hub)
		game.endPicking(hub)
	}
}

func (game *Game) getPickableIntCardIds(playerId string) []uint16 {
	var cardIds []uint16
	remainingCardIds := game.getCurrentRound().getRemainingCardIds(playerId)

	var trick = game.getCurrentTrick()

	table := newTable(
		trick.getAllPickedCardIds(),
	)

	hand := newHand(remainingCardIds)
	pickableCardIds := hand.pickables(table)

	for _, pickableCardId := range pickableCardIds {
		cardIds = append(cardIds, uint16(pickableCardId))
	}

	return cardIds
}

func (game *Game) endPicking(hub *Hub) {
	game.ExpirationTime = 0

	if game.isTrickOver() {
		game.announceTrickWinner(hub)
		game.nextTrick(hub)
	} else {
		game.startPicking(hub)
	}
}

func (game *Game) isTrickOver() bool {
	var trick = game.getCurrentTrick()
	return len(trick.PickedCards) == len(game.Players)
}

func (game *Game) announceTrickWinner(hub *Hub) {
	cardId, playerId := game.findTrickWinner()

	game.getCurrentTrick().WinnerPlayerId = playerId
	game.getCurrentTrick().WinnerCardId = cardId

	content := responses.AnnounceTrickWinner{
		PlayerId: playerId,
		CardId:   uint16(cardId),
	}

	m := &ServerMessage{
		Content: content,
		Command: constants.CommandAnnounceTrickWinner,
		GameId:  game.Id,
	}

	hub.Dispatch <- m
}

func (game *Game) announceScores(hub *Hub) {
	var round = game.getCurrentRound()
	round.calculateScores()

	content := responses.AnnounceScore{}

	for playerId, score := range round.Scores {
		game.Players[playerId].Score += score
		s := responses.Score{
			PlayerId: playerId,
			Score:    game.Players[playerId].Score,
		}
		content.Scores = append(content.Scores, s)
	}

	m := &ServerMessage{
		Content: content,
		Command: constants.CommandAnnounceScores,
		GameId:  game.Id,
	}

	hub.Dispatch <- m
}

func (game *Game) nextTrick(hub *Hub) {
	if game.Trick == game.Round {
		game.announceScores(hub)
		game.NextRound(hub)
		return
	}

	game.Trick++
	game.getCurrentRound().Tricks[game.Trick-1] = &Trick{
		Number: game.Trick,
	}

	content := responses.NextTrick{
		Round: game.Round,
		Trick: game.Trick,
	}

	m := &ServerMessage{
		Content: content,
		Command: constants.CommandNextTrick,
		GameId:  game.Id,
	}

	hub.Dispatch <- m

	game.startPicking(hub)
}

func (game *Game) findTrickWinner() (CardId, string) {
	var trick = game.getCurrentTrick()

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

func (game *Game) pickForIdlePlayer(hub *Hub) {
	var trick = game.getCurrentTrick()
	var pickerId = trick.PickingUserId

	if trick.isPlayerPicked(pickerId) {
		return
	}

	availableCardIds := game.getPickableIntCardIds(pickerId)

	pickedCard := PickedCard{
		PlayerId: pickerId,
		CardId:   CardId(availableCardIds[0]),
	}
	trick.PickedCards = append(trick.PickedCards, pickedCard)

	content := responses.Picked{
		PlayerId: pickerId,
		CardId:   uint16(pickedCard.CardId),
	}

	m := &ServerMessage{
		Content: content,
		Command: constants.CommandPicked,
		GameId:  game.Id,
	}

	hub.Dispatch <- m
}

func (game *Game) Initialize(hub *Hub, receiverId string) {
	var players []responses.Player

	for _, player := range game.getPlayers() {
		var p responses.Player

		p.Id = player.Id
		p.Username = player.Username
		p.Avatar = player.Avatar
		p.Score = player.Score

		if game.Round != 0 {
			var round = game.getCurrentRound()

			p.WonTricksCount = round.getWonTricksCount(player.Id)

			if player.Id == receiverId || game.State != constants.StateBidding {
				p.Bid = round.Bids[player.Id]
			}

			// Receiver must not be aware of other cards
			if player.Id == receiverId {
				p.HandCardIds = round.getRemainingIntCardIds(player.Id)
				p.PickableCardIds = game.getPickableIntCardIds(receiverId)
			}
		}

		players = append(players, p)
	}

	content := responses.Init{
		Round:          game.Round,
		Trick:          game.Trick,
		State:          game.State,
		ExpirationTime: game.ExpirationTime,
		Players:        players,
		CreatorId:      game.CreatorId,
	}

	if game.Round != 0 {
		var trick = game.getCurrentTrick()
		content.PickingUserId = trick.PickingUserId
		content.TableCards = trick.getAllPickedCards()
	}

	m := &ServerMessage{
		Content:    content,
		Command:    constants.CommandInit,
		GameId:     game.Id,
		ReceiverId: receiverId,
	}

	hub.Dispatch <- m
}

func (game *Game) validateUserPickedCard(pickedCardId uint16, playerId string) error {
	if game.State != constants.StatePicking {
		return errors.New("we are not accepting picking command in this state")
	}

	var trick = game.getCurrentTrick()

	if trick.PickingUserId != playerId {
		return errors.New("it's not your turn to pick a card")
	}

	cardIds := game.getPickableIntCardIds(playerId)

	var exists = false
	for _, cardId := range cardIds {
		if cardId == pickedCardId {
			exists = true
		}
	}

	if !exists {
		return errors.New("you don't own the card")
	}

	return nil
}

func (game *Game) Pick(hub *Hub, cardId uint16, playerId string) {

	err := game.validateUserPickedCard(cardId, playerId)

	if err != nil {
		//
		content := responses.Error{
			Message: err.Error(),
			StatusCode: 422,
		}
		m := &ServerMessage{
			Content:    content,
			GameId:     game.Id,
			ReceiverId: playerId,
		}
		hub.Dispatch <- m
		return
	}

	pickedCard := PickedCard{
		PlayerId: playerId,
		CardId:   CardId(cardId),
	}

	var trick = game.getCurrentTrick()

	trick.PickedCards = append(trick.PickedCards, pickedCard)

	content := responses.Picked{
		PlayerId: playerId,
		CardId:   cardId,
	}

	var m = &ServerMessage{
		Content: content,
		Command: constants.CommandPicked,
		GameId:  game.Id,
	}

	hub.Dispatch <- m

	game.endPicking(hub)
}

func (game *Game) GetAvailableAvatar() string {
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

func (game *Game) Bid(hub *Hub, playerId string, number int) {
	if number < 0 || number > game.Round {
		content := responses.Error{
			Message: "Invalid bid number.",
			StatusCode: 422,
		}
		m := &ServerMessage{
			Content:    content,
			GameId:     game.Id,
			ReceiverId: playerId,
		}
		hub.Dispatch <- m
		return
	}
	game.getCurrentRound().Bids[playerId] = number
	content := responses.Bade{Number: number}
	m := &ServerMessage{
		Content:    content,
		Command:    constants.CommandBade,
		GameId:     game.Id,
		ReceiverId: playerId,
	}
	hub.Dispatch <- m
	// TODO: If he is the last one bidding, call endBidding()
}

func (game *Game) Join(hub *Hub, player *Player) {
	m := &ServerMessage{
		Command: constants.CommandJoined,
		Content: responses.Player{
			Id:       player.Id,
			Username: player.Username,
			Avatar:   player.Avatar,
		},
		GameId:     player.GameId,
		ExcludedId: player.Id,
	}

	hub.Dispatch <- m
}

func (game *Game) Left(hub *Hub, playerId string) {
	m := &ServerMessage{
		Content: responses.Left{PlayerId: playerId},
		Command: constants.CommandLeft,
		GameId:  game.Id,
	}

	hub.Dispatch <- m
}

func (game *Game) getBiddingExpirationDuration() time.Duration {
	// As the round number increases, it takes more time to complete the card dealing animation.
	// Therefore, we need to increase the wait time for each level
	// Each animation takes about 2 seconds
	return constants.WaitTime + time.Duration(game.Round)*2*time.Second
}

func (game *Game) getCurrentRound() *Round {
	return game.Rounds[game.Round-1]
}

func (game *Game) getPreviousRound() *Round {
	return game.Rounds[game.Round-2]
}

func (game *Game) getCurrentTrick() *Trick {
	var round = game.getCurrentRound()
	return round.Tricks[game.Trick-1]
}

func (game *Game) getPreviousTrick() *Trick {
	return game.getCurrentRound().Tricks[game.Trick-2]
}

func (game *Game) getPlayers() []*Player {
	var players = make([]*Player, 0, len(game.Players))

	var playerIds = make([]string, 0, len(game.Players))

	for playerId := range game.Players {
		playerIds = append(playerIds, playerId)
	}

	sort.SliceStable(playerIds, func(i, j int) bool {
		return game.Players[playerIds[i]].Index < game.Players[playerIds[j]].Index
	})

	for _, playerId := range playerIds {
		players = append(players, game.Players[playerId])
	}

	return players
}
