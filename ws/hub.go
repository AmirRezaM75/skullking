package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Room struct {
	Id      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
	Stall      chan bool
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message),
		Stall:      make(chan bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case expired, ok := <-h.Stall:
			fmt.Println("Stall", expired, ok)
			if !ok {
				h.Stall = nil
			} else {
				message := &Message{
					Content: "Expired",
					RoomId:  "xxx-yyy-zzz",
				}
				if _, ok := h.Rooms[message.RoomId]; ok {
					for _, client := range h.Rooms[message.RoomId].Clients {
						client.Message <- message
					}
				}
			}
		case client := <-h.Register:
			fmt.Println("Register")
			if _, ok := h.Rooms[client.RoomId]; ok {
				room := h.Rooms[client.RoomId]
				if _, ok = room.Clients[client.Id]; !ok {
					room.Clients[client.Id] = client
				}
			}
		case client := <-h.Unregister:
			fmt.Println("Unregister")
			if _, ok := h.Rooms[client.RoomId]; ok {
				clients := h.Rooms[client.RoomId].Clients
				if _, ok = clients[client.Id]; ok {
					if len(h.Rooms[client.RoomId].Clients) != 0 {
						h.Broadcast <- &Message{
							Content: fmt.Sprintf("%s left the room.", client.Id),
							RoomId:  client.RoomId,
						}
					}
					delete(clients, client.Id)
					close(client.Message)
				}
			}
		case message := <-h.Broadcast:
			fmt.Println("Broadcast")

			if _, ok := h.Rooms[message.RoomId]; ok {
				for _, client := range h.Rooms[message.RoomId].Clients {
					client.Message <- message
				}
			}

			if message.Content == "start" {
				var deck Deck
				cards := generateCards()
				deck.cards = cards
				deck.shuffle()
				items := deck.deal(len(h.Rooms[message.RoomId].Clients), 3)
				if _, ok := h.Rooms[message.RoomId]; ok {
					index := 0
					for _, client := range h.Rooms[message.RoomId].Clients {
						userCards, _ := json.Marshal(items[index])
						index++
						cardsMessage := &Message{
							ContentType: "json",
							Content:     string(userCards),
							Command:     "DEAL_CARDS",
							RoomId:      client.RoomId,
						}
						client.Message <- cardsMessage
					}
				}
			}
		}
	}
}

func wait(hub *Hub) {
	log.Println("wait method.")
	time.Sleep(time.Second * 5)
	hub.Stall <- true
	close(hub.Stall)
}
