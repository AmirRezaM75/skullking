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
			if !ok {
				h.Stall = nil
			} else {
				if expired {
					message := &Message{
						Command: CommandBettingEnded,
						RoomId:  "xxx-yyy-zzz",
					}
					if _, ok := h.Rooms[message.RoomId]; ok {
						for _, client := range h.Rooms[message.RoomId].Clients {
							client.Message <- message
						}
					}
				}
			}
		case client := <-h.Register:
			if _, ok := h.Rooms[client.RoomId]; ok {
				room := h.Rooms[client.RoomId]
				if _, ok = room.Clients[client.Id]; !ok {
					room.Clients[client.Id] = client
				}
			}
		case client := <-h.Unregister:
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

					betCommand := BetCommand{Round: 3, EndsAt: time.Now().Add(WaitTime).Unix()}
					betCommandJson, _ := json.Marshal(betCommand)
					for _, client := range h.Rooms[message.RoomId].Clients {
						cardsMessage := &Message{
							ContentType: "json",
							Content:     string(betCommandJson),
							Command:     CommandBettingStarted,
							RoomId:      client.RoomId,
						}
						client.Message <- cardsMessage
					}
				}
				go wait(h)
			}
		}
	}
}

func wait(hub *Hub) {
	log.Println("wait method.")
	time.Sleep(WaitTime)
	hub.Stall <- true
	close(hub.Stall)
}

const WaitTime = 10 * time.Second

// Commands

const CommandBettingStarted = "BETTING_STARTED"
const CommandBettingEnded = "BETTING_ENDED"

// Command Structs

type BetCommand struct {
	Round  int32 `json:"round"`
	EndsAt int64 `json:"endsAt"`
}
