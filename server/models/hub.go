package models

import (
	"skullking/constants"
	"skullking/pkg/syncx"
	"time"
)

// Couldn't define this interface inside contracts because of 'import cycle not allowed' error

type GameRepository interface {
	Create(u *Game) error
}

type PublisherService interface {
	Publish(message string) error
}

type LogService interface {
	Error(map[string]string)
	Info(map[string]string)
}

type Hub struct {
	Games            syncx.Map[string, *Game]
	Dispatch         chan *ServerMessage
	GameRepository   GameRepository
	PublisherService PublisherService
	LogService       LogService
}

func NewHub(
	gameRepository GameRepository,
	publisherService PublisherService,
	logService LogService,
) *Hub {
	return &Hub{
		Games:            syncx.Map[string, *Game]{},
		Dispatch:         make(chan *ServerMessage),
		GameRepository:   gameRepository,
		PublisherService: publisherService,
		LogService:       logService,
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

func (h *Hub) Unsubscribe(player *Player) {
	// If the game has already started, we will not remove the player,
	// and the server decide on behalf of them.
	if game, ok := h.Games.Load(player.GameId); ok {
		game.Left(h, player.Id)
	}

	// TODO: If every one left the game delete the game.
}

func (h *Hub) Cleanup() {
	h.Games.Range(func(_ string, game *Game) bool {
		if game.CreatedAt <= time.Now().Add(-30*time.Minute).Unix() &&
			game.State == constants.StatePending {
			h.LogService.Info(map[string]string{
				"message": "Delete game due to inactivity",
				"gameId":  game.Id,
			})
			h.Games.Delete(game.Id)
		}
		return true
	})
}
