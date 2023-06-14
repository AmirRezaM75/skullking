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
}

func (player *Player) Write() {
	defer func() {
		_ = player.Connection.Close()
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
		hub.Unsubscribe(player)
		_ = player.Connection.Close()
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
		game.CreatorId == player.Id {
		game.NextRound(hub)
	}

	if message.Command == constants.CommandPick {
		cardId, _ := strconv.Atoi(message.Content)
		game.Pick(hub, cardId, player.Id)
	}
}
