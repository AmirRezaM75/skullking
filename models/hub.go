package models

import (
	"encoding/json"
	"skullking/constants"
	"skullking/pkg/syncx"
	"time"
)

// Couldn't define this interface inside contracts because of 'import cycle not allowed' error

type GameRepository interface {
	Create(u *Game) error
}

type BotRepository interface {
	Bid(cardIds []uint16) (int, error)
	Pick(
		handCards []uint16,
		pickableCards []uint16,
		tableCards []uint16,
		observedCards []uint16,
		bid int,
		tricksTaken uint,
		playerIndex int,
		numPlayers int,
	) (uint16, error)
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
	BotRepository    BotRepository
	PublisherService PublisherService
	LogService       LogService
}

func NewHub(
	gameRepository GameRepository,
	botRepository BotRepository,
	publisherService PublisherService,
	logService LogService,
) *Hub {
	return &Hub{
		Games:            syncx.Map[string, *Game]{},
		Dispatch:         make(chan *ServerMessage),
		GameRepository:   gameRepository,
		BotRepository:    botRepository,
		PublisherService: publisherService,
		LogService:       logService,
	}
}

func (h *Hub) logMessage(message *ServerMessage) {
	content, err := json.Marshal(message)

	if err != nil {
		h.LogService.Error(map[string]string{
			"message":     err.Error(),
			"description": "Can not marshal message",
		})
		content = []byte("")
	}

	h.LogService.Info(map[string]string{
		"message":    "event dispatched",
		"excludedId": message.ExcludedId,
		"receiverId": message.ReceiverId,
		"gameId":     message.GameId,
		"command":    message.Command,
		"content":    string(content),
	})
}

func (h *Hub) Run() {
	for {
		select {
		case message := <-h.Dispatch:
			if _, ok := h.Games.Load(message.GameId); ok {
				var game, _ = h.Games.Load(message.GameId)

				h.logMessage(message)

				// If there is no specific receiver broadcast it to all players
				if message.ReceiverId == "" {
					game.Players.Range(func(_ string, player *Player) bool {
						player.mutex.Lock()
						defer player.mutex.Unlock()

						if message.ExcludedId != player.Id && !player.IsClosed {
							player.Message <- message
						}

						return true
					})
				} else {
					if player, ok := game.Players.Load(message.ReceiverId); ok {
						func() {
							player.mutex.Lock()
							defer player.mutex.Unlock()
							if !player.IsClosed {
								player.Message <- message
							}
						}()
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
			(game.State == constants.StatePending ||
				game.State == constants.StateEnded) {
			// Quick way to fix nil pointer when running tests
			if h.LogService != nil {
				h.LogService.Info(map[string]string{
					"message": "Game has been deleted.",
					"gameId":  game.Id,
					"state":   game.State,
				})
			}
			h.Games.Delete(game.Id)
		}
		return true
	})
}
