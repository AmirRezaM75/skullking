package models

import (
	"fmt"
	"github.com/AmirRezaM75/skull-king/constants"
	"github.com/AmirRezaM75/skull-king/pkg/syncx"
	"time"
)

// TODO: Couldn't define this interface inside contracts because of 'import cycle not allowed' error

type GameRepository interface {
	Create(u *Game) error
}

type Hub struct {
	Games          syncx.Map[string, *Game]
	Dispatch       chan *ServerMessage
	GameRepository GameRepository
}

func NewHub(gameRepository GameRepository) *Hub {
	return &Hub{
		Games:          syncx.Map[string, *Game]{},
		Dispatch:       make(chan *ServerMessage),
		GameRepository: gameRepository,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case message := <-h.Dispatch:
			if _, ok := h.Games.Load(message.GameId); ok {
				var game, _ = h.Games.Load(message.GameId)
				// If there is no specific receiver broadcast it to all players
				if message.ReceiverId == "" {
					game.Players.Range(func(_ string, player *Player) bool {
						if message.ExcludedId != player.Id && player.IsConnected {
							player.Message <- message
						}
						return true
					})
				} else {
					if p, ok := game.Players.Load(message.ReceiverId); ok && p.IsConnected {
						p.Message <- message
					}
				}
			}
		}
	}
}

func (h *Hub) Subscribe(player *Player) {
	if game, ok := h.Games.Load(player.GameId); ok {
		game.Players.Store(player.Id, player)
	}
}

func (h *Hub) Unsubscribe(player *Player) {
	// If the game status is PENDING, we will remove the player from the game
	// to inform the game creator of the total number of players before starting.
	// However, if the game has already started, we will not remove the player,
	// and the server decide on behalf of them.
	if game, ok := h.Games.Load(player.GameId); ok {
		if game.State == constants.StatePending {
			game.Left(h, player.Id)
			game.Players.Delete(player.Id)
		}
	}

	// TODO: If every one left the game delete the game.
}

func (h *Hub) Cleanup() {
	h.Games.Range(func(_ string, game *Game) bool {
		if game.CreatedAt <= time.Now().Add(-30*time.Minute).Unix() &&
			game.State == constants.StatePending {
			fmt.Printf("Delete game %s due to inactivity.\n", game.Id)
			h.Games.Delete(game.Id)
		}
		return true
	})
}
