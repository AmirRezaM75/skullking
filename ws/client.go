package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type Client struct {
	Connection       *websocket.Conn
	Message          chan *Message
	Bid              int    `json:"bid"`
	Id               string `json:"id"`
	RoomId           string `json:"roomId"`
	LastPickedCardId int
	CardIds          []int
}

// Message from server to client
type Message struct {
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
	Command     string `json:"command"`
	RoomId      string `json:"roomId"`
	ReceiverId  string
	SenderId    string
}

func (c *Client) writeMessage() {
	defer func() {
		c.Connection.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}

		c.Connection.WriteJSON(message)
	}
}

type MessageStruct struct {
	Command string `json:"command"`
	Content string `json:"content"`
}

func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Connection.Close()
	}()

	var message MessageStruct

	for {
		_, m, err := c.Connection.ReadMessage()
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

		msg := &Message{
			ContentType: "string",
			Content:     message.Content,
			Command:     message.Command,
			RoomId:      c.RoomId,
			ReceiverId:  "",
			SenderId:    c.Id,
		}
		fmt.Printf("%+v\n", msg)

		if msg.Command == CommandBet && hub.Rooms[c.RoomId].Status == StatusMakingBid {
			c.Bid, _ = strconv.Atoi(msg.Content)
			continue
		}

		if msg.Command == CommandPick {
			log.Println("LastPickingUserId", hub.Rooms[c.RoomId].LastPickingUserId)
			log.Println("ClientId", c.Id)
			if hub.Rooms[c.RoomId].LastPickingUserId == c.Id {
				// TODO: Check if cardId is valid and exists in the last dealt cards
				cardId, _ := strconv.Atoi(msg.Content)
				c.LastPickedCardId = cardId
			} else {
				continue
			}
		}

		hub.Broadcast <- msg
	}
}

type Card struct {
	Id     int    `json:"id"`
	Color  string `json:"color"`
	Number int    `json:"number"`
	Group  string `json:"group"`
}

func generateCards() []Card {
	colors := [4]string{"Black", "Red", "Green", "Yellow"}

	var cards []Card
	id := 1
	for _, color := range colors {
		for i := 1; i <= 14; i++ {
			card := Card{
				Id:     id,
				Color:  color,
				Number: i,
				Group:  "parrot",
			}
			cards = append(cards, card)
			id++
		}
	}

	for i := 1; i <= 5; i++ {
		card := Card{
			Id:     id,
			Color:  "White",
			Number: 0,
			Group:  "escape",
		}
		cards = append(cards, card)
	}

	for i := 1; i <= 5; i++ {
		card := Card{
			Id:     id,
			Color:  "Brown",
			Number: 0,
			Group:  "Pirate",
		}
		cards = append(cards, card)
	}

	for i := 1; i <= 2; i++ {
		card := Card{
			Id:     id,
			Color:  "pink",
			Number: 0,
			Group:  "mermaid",
		}
		cards = append(cards, card)
	}

	card := Card{
		Id:     id,
		Color:  "Black",
		Number: 0,
		Group:  "skullKing",
	}
	cards = append(cards, card)

	card = Card{
		Id:     id,
		Color:  "White",
		Number: 0,
		Group:  "Kraken",
	}
	cards = append(cards, card)

	return cards
}

type Deck struct {
	cards []Card
}

func (d Deck) shuffle() {
	cards := d.cards

	for i := range cards {
		source := rand.NewSource(time.Now().UnixNano())
		r := rand.New(source)
		j := r.Intn(len(d.cards) - 1)
		cards[i], cards[j] = cards[j], cards[i]
	}
}

func (d Deck) deal(count, size int) [][]Card {
	var output [][]Card

	index := 0

	for i := 0; i < count; i++ {
		output = append(output, d.cards[index:size+index])
		index = size + index
	}

	return output
}
