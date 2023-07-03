package models

import (
	"encoding/json"
	"github.com/AmirRezaM75/skull-king/constants"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
)

type Player struct {
	Id         string
	Username   string
	GameId     string
	Avatar     string
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
		// to avoid close of closed channel panic
		// first we need to check if channel is open or not
		if _, ok := <-player.Message; ok {
			close(player.Message)
		}
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
		hub.Unsubscribe(player)
	}()

	var message ClientMessage

	for {
		_, m, err := player.Connection.ReadMessage()
		//TODO: Validation of command types.
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		err = json.Unmarshal(m, &message)

		if err != nil {
			log.Printf("unmarshal error: %v", err)
			continue
		}

		player.react(message, hub)
	}
}

func (player *Player) react(message ClientMessage, hub *Hub) {
	var game = hub.Games[player.GameId]

	if message.Command == constants.CommandBid && game.State == constants.StateBidding {
		number, err := strconv.Atoi(message.Content)

		if err != nil {
			number = 0
		}

		game.Bid(hub, player.Id, number)
		return
	}

	if message.Command == constants.CommandStart &&
		game.State == constants.StatePending &&
		game.CreatorId == player.Id &&
		len(game.Players) > 1 {
		game.NextRound(hub)
	}

	if message.Command == constants.CommandPick {
		cardId, _ := strconv.Atoi(message.Content)
		game.Pick(hub, uint16(cardId), player.Id)
	}
}
