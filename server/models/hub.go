package models

import (
	"fmt"
	"github.com/AmirRezaM75/skull-king/constants"
	"time"
)

// TODO: Couldn't define this interface inside contracts because of 'import cycle not allowed' error
type GameRepository interface {
	Create(u Game) (*Game, error)
}

type Hub struct {
	Games          map[string]*Game
	Dispatch       chan *ServerMessage
	GameRepository GameRepository
}

func NewHub(gameRepository GameRepository) *Hub {
	return &Hub{
		Games:          make(map[string]*Game),
		Dispatch:       make(chan *ServerMessage),
		GameRepository: gameRepository,
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
	// and the server decide on behalf of them.
	if _, ok := h.Games[player.GameId]; ok {
		game := h.Games[player.GameId]
		if game.State == constants.StatePending {
			game.Left(h, player.Id)
			delete(game.Players, player.Id)
		}
	}
	// TODO: If every one left the game delete the game.
}

func (h *Hub) Cleanup() {
	for _, game := range h.Games {
		if game.CreatedAt <= time.Now().Add(-30*time.Minute).Unix() &&
			game.State == constants.StatePending {
			fmt.Printf("Delete game %s due to inactivity.\n", game.Id)
			delete(h.Games, game.Id)
		}
	}
}
