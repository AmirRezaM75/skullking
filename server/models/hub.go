package models

import (
	"github.com/AmirRezaM75/skull-king/constants"
	"time"
)

type Hub struct {
	Games      map[string]*Game
	Register   chan *Player
	Unregister chan *Player
	Dispatch   chan *ServerMessage
	Remind     chan Reminder
}

type Reminder struct {
	gameId string
}

func NewHub() *Hub {
	return &Hub{
		Games:      make(map[string]*Game),
		Register:   make(chan *Player),
		Unregister: make(chan *Player),
		Dispatch:   make(chan *ServerMessage, 10),
		Remind:     make(chan Reminder),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case reminder, ok := <-h.Remind:
			game := h.Games[reminder.gameId]
			if !ok {
				h.Remind = nil
			} else {
				if game.State == constants.StatePicking {
					game.EndPicking(h)
				}
			}
		case player := <-h.Register:
			if _, ok := h.Games[player.GameId]; ok {
				h.Games[player.GameId].Players[player.Id] = player
			}
		case player := <-h.Unregister:
			// Remove player from the game when game status is PENDING
			// Because we need to inform game creator how many people
			// are in the game before starting.
			if _, ok := h.Games[player.GameId]; ok {
				game := h.Games[player.GameId]
				if game.State == constants.StatePending {
					game.Left(h, player.Id)
					delete(game.Players, player.Id)
				}
			}
			// Not going to remove player from game
			// because it should continue playing till the end
			// close(player.ServerMessage)
			// If every one left the game delete the game.

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

func (h *Hub) newReminder(gameId string) {
	time.Sleep(constants.WaitTime)

	h.Remind <- Reminder{
		gameId: gameId,
	}

	//close(hub.remind)
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
