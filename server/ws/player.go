package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
)

type Player struct {
	id         int
	gameId     string
	avatar     string
	score      int
	connection *websocket.Conn
	message    chan *ServerMessage
}

// ServerMessage Message from server to client structure
type ServerMessage struct {
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
	Command     string `json:"command"`
	SenderId    int    `json:"senderId"`
	gameId      string
	receiverId  int
}

// ClientMessage Message from client to server structure
type ClientMessage struct {
	Command string `json:"command"`
	Content string `json:"content"`
}

func (p *Player) writeMessage() {
	defer func() {
		_ = p.connection.Close()
	}()

	for {
		message, ok := <-p.message

		if !ok {
			return
		}

		_ = p.connection.WriteJSON(message)
	}
}

func (p *Player) readMessage(hub *Hub) {
	defer func() {
		hub.unregister <- p
		_ = p.connection.Close()
	}()

	var message ClientMessage

	for {
		_, m, err := p.connection.ReadMessage()
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

		p.react(message, hub)
	}
}

func (p *Player) react(message ClientMessage, hub *Hub) {
	var game = hub.games[p.gameId]

	if message.Command == CommandInit {
		initialize(game, hub, p.id)
	}

	if message.Command == CommandBid && game.state == StateBidding {
		game.rounds[game.round].bids[p.id], _ = strconv.Atoi(message.Content)
		// TODO: If he is the last one picking the card, clear the timer
		return
	}

	if message.Command == CommandStart && game.state == "" {
		start(game, hub)
	}

	if message.Command == CommandPick {
		cardId, _ := strconv.Atoi(message.Content)
		pick(game, hub, cardId, p.id)
	}
}
