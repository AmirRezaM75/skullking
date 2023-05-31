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

	if message.Command == constants.CommandInit {
		game.Initialize(hub, player.Id)
	}

	if message.Command == constants.CommandBid && game.State == constants.StateBidding {
		game.Rounds[game.Round].Bids[player.Id], _ = strconv.Atoi(message.Content)
		// TODO: If he is the last one picking the card, clear the timer
		return
	}

	if message.Command == constants.CommandStart && game.State == "" {
		game.Start(hub)
	}

	if message.Command == constants.CommandPick {
		cardId, _ := strconv.Atoi(message.Content)
		game.Pick(hub, cardId, player.Id)
	}
}
