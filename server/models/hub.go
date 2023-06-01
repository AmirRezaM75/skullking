package models

import (
	"github.com/AmirRezaM75/skull-king/constants"
)

type Hub struct {
	Games    map[string]*Game
	Dispatch chan *ServerMessage
}

type Reminder struct {
	gameId string
}

func NewHub() *Hub {
	return &Hub{
		Games:    make(map[string]*Game),
		Dispatch: make(chan *ServerMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case message := <-h.Dispatch:
			// If there is no specific receiver broadcast it to all players
			if _, ok := h.Games[message.GameId]; ok {
				var game = h.Games[message.GameId]
				if message.ReceiverId == "" {
					for _, player := range game.Players {
						player.Message <- message
					}
				} else {
					if _, ok = game.Players[message.ReceiverId]; ok {
						game.Players[message.ReceiverId].Message <- message
					}
				}
			}
		}
	}
}

func (h *Hub) Subscribe(player *Player) {
	if _, ok := h.Games[player.GameId]; ok {
		h.Games[player.GameId].Players[player.Id] = player
	}
}

func (h *Hub) Unsubscribe(player *Player) {
	// If the game status is PENDING, we will remove the player from the game
	// to inform the game creator of the total number of players before starting.
	// However, if the game has already started, we will not remove the player,
	// and the server decide on behalf of them
	if _, ok := h.Games[player.GameId]; ok {
		game := h.Games[player.GameId]
		if game.State == constants.StatePending {
			game.Left(h, player.Id)
			delete(game.Players, player.Id)
		}
	}
	// TODO: If every one left the game delete the game.
}

type UserBet struct {
	UserId string `json:"userId"`
	Bet    int    `json:"bet"`
}

type UserPickedCard struct {
	UserId string `json:"userId"`
	CardId int    `json:"cardId"`
}

type CommandPickContent struct {
	CardId int `json:"cardId"`
}
