package models

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"skullking/constants"
	"strconv"
	"sync"
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
	// To keep track of closed channel
	IsClosed bool
	mutex    sync.Mutex
}

// Different scenarios for 'close of closed channel'
// 1) If user opens duplicate tab and close the first one

func (player *Player) Kick() {

	// We are using mutex to make sure IsClosed value is evaluated correctly
	// when reading its value at the same time.
	// https://go101.org/article/channel-closing.html
	player.mutex.Lock()

	defer player.mutex.Unlock()

	if !player.IsClosed {
		close(player.Message)
		player.IsClosed = true
	}

	// First we need to check if it's nil or not
	// we call kick method in game_handler, and player may has no connection
	if player.Connection != nil {
		_ = player.Connection.Close()
	}

	player.IsConnected = false
}

func (player *Player) Write() {
	defer func() {
		player.Kick()
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
		player.Kick()
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
