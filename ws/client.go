package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"time"
)

type Client struct {
	Connection *websocket.Conn
	Message    chan *Message
	Id         string `json:"id"`
	RoomId     string `json:"roomId"`
}

type Message struct {
	Content string `json:"content"`
	RoomId  string `json:"roomId"`
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

func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Connection.Close()
	}()

	for {
		_, m, err := c.Connection.ReadMessage()
		fmt.Println(m)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := &Message{
			Content: string(m),
			RoomId:  c.RoomId,
		}

		if msg.Content == "start" {
			fmt.Println("Starting...")
			go wait(hub)
		}

		hub.Broadcast <- msg
	}
}

type Card struct {
	Color  string `json:"color"`
	Number int    `json:"number"`
}

func generateCards() []Card {
	colors := [4]string{"Black", "Red", "Green", "Yellow"}

	var cards []Card

	for _, color := range colors {
		for i := 1; i <= 14; i++ {
			card := Card{
				Color:  color,
				Number: i,
			}
			cards = append(cards, card)
		}
	}

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
