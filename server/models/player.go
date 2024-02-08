package models

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"skullking/constants"
	"strconv"
)

type Player struct {
	Id         string
	Username   string
	GameId     string
	AvatarId   uint8
	Score      int
	Index      int
	Connection *websocket.Conn
	Message    chan *ServerMessage
	// To determine whether a player is currently participating in the game or has left,
	// it can be useful to check whether their corresponding channel is open or closed.
	IsConnected bool
}

func (player *Player) disconnect() {
	if player.IsConnected {
		_ = player.Connection.Close()
		player.IsConnected = false
	}
}

func (player *Player) Write() {
	defer func() {
		player.disconnect()
	}()

	for {
		message, ok := <-player.Message

		if !ok {
			return
		}

		_ = player.Connection.WriteJSON(message)
	}
}

func (player *Player) Read(hub *Hub) {
	defer func() {
		player.disconnect()
		close(player.Message)
		hub.Unsubscribe(player)
	}()

	var message ClientMessage

	for {
		_, m, err := player.Connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		err = json.Unmarshal(m, &message)

		if err != nil {
			hub.LogService.Error(map[string]string{
				"message":     err.Error(),
				"method":      "player@read",
				"description": "Could not unmarshal message.",
			})
			continue
		}

		player.react(message, hub)
	}
}

func (player *Player) react(message ClientMessage, hub *Hub) {
	if _, ok := hub.Games.Load(player.GameId); !ok {
		return
	}

	game, _ := hub.Games.Load(player.GameId)

	if message.Command == constants.CommandBid &&
		game.State == constants.StateBidding {
		number, err := strconv.Atoi(message.Content)

		if err != nil {
			number = 0
		}

		game.Bid(hub, player.Id, number)
		return
	}

	if message.Command == constants.CommandPick {
		cardId, _ := strconv.Atoi(message.Content)
		game.Pick(hub, uint16(cardId), player.Id)
	}
}
