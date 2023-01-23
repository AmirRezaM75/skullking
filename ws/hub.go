package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Room struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	Clients           map[string]*Client
	Round             int
	Status            string
	LastPickingUserId string
}

const StatusNotStarted = "NOT_STARTED"
const StatusMakingBid = "MAKING_BID"
const StatusEndOfAuction = "END_OF_AUCTION"
const StatusTakingTrick = "TAKING_A_TRICK"

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
			fmt.Println("after wait.", expired, ok)
			if !ok {
				h.Stall = nil
			} else {

				var userBets []UserBet
				if expired {

					if _, ok := h.Rooms["xxx-yyy-zzz"]; ok {
						for _, client := range h.Rooms["xxx-yyy-zzz"].Clients {
							userBets = append(userBets, UserBet{
								UserId: client.Id,
								Bet:    client.Bid,
							})
						}
					}

					content, _ := json.Marshal(userBets)
					message := &Message{
						Command:     CommandBettingEnded,
						Content:     string(content),
						ContentType: "json",
						RoomId:      "xxx-yyy-zzz",
					}
					if _, ok := h.Rooms[message.RoomId]; ok {
						for _, client := range h.Rooms[message.RoomId].Clients {
							client.Message <- message
						}
					}
					h.Rooms[message.RoomId].Status = StatusEndOfAuction

					game := h.Rooms[message.RoomId]

					// Get first client
					for userId, _ := range game.Clients {
						pickingStartedContent := PickingStartedContent{
							UserId: userId,
						}
						content, _ = json.Marshal(pickingStartedContent)
						message = &Message{
							Command:     CommandPickingStarted,
							ContentType: "json",
							Content:     string(content),
							// No access to roomId
							RoomId: "xxx-yyy-zzz",
						}
						game.LastPickingUserId = userId
						break
					}

					if _, ok = h.Rooms[message.RoomId]; ok {
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
			if message.Command == CommandInitGame {
				if _, ok := h.Rooms[message.RoomId]; ok {
					if _, ok := h.Rooms[message.RoomId].Clients[message.ReceiverId]; ok {
						h.Rooms[message.RoomId].Clients[message.ReceiverId].Message <- message
					}
				}
			}

			// Broadcast to everyone except the sender
			if message.Command == CommandUserJoined {
				if _, ok := h.Rooms[message.RoomId]; ok {
					for _, client := range h.Rooms[message.RoomId].Clients {
						if client.Id == message.SenderId {
							continue
						}
						client.Message <- message
					}
				}
			}

			if message.Command == CommandStart {
				if _, ok := h.Rooms[message.RoomId]; ok {
					room := h.Rooms[message.RoomId]
					room.Round++
					var deck Deck
					cards := generateCards()
					deck.cards = cards
					deck.shuffle()
					items := deck.deal(len(room.Clients), 3)
					index := 0
					for _, client := range room.Clients {
						userCards, _ := json.Marshal(items[index])
						index++
						cardsMessage := &Message{
							ContentType: "json",
							Content:     string(userCards),
							Command:     CommandDealCards,
							RoomId:      client.RoomId,
						}
						client.Message <- cardsMessage
					}

					betCommand := BetCommand{
						Round:  room.Round,
						EndsAt: time.Now().Add(WaitTime).Unix(),
					}
					betCommandJson, _ := json.Marshal(betCommand)
					for _, client := range room.Clients {
						cardsMessage := &Message{
							ContentType: "json",
							Content:     string(betCommandJson),
							Command:     CommandBettingStarted,
							RoomId:      client.RoomId,
						}
						client.Message <- cardsMessage
					}
					room.Status = StatusMakingBid
				}
				go wait(h)
			}

			// Broadcast to everyone
			if message.Command == CommandPick {
				if _, ok := h.Rooms[message.RoomId]; ok {
					for _, client := range h.Rooms[message.RoomId].Clients {
						client.Message <- message
					}
				}
			}

			if message.Command == CommandPick {
				game := h.Rooms[message.RoomId]
				nextUserFound := false
				var nextUserId string
				for userId, _ := range game.Clients {
					if nextUserFound {
						nextUserId = userId
					}
					if userId == game.LastPickingUserId {
						nextUserFound = true
					}
				}
				pickingStartedContent := PickingStartedContent{
					UserId: nextUserId,
					EndsAt: time.Now().Add(WaitTime).Unix(),
				}
				content, _ := json.Marshal(pickingStartedContent)
				message = &Message{
					Command:     CommandPickingStarted,
					ContentType: "json",
					Content:     string(content),
					// No access to roomId
					RoomId: "xxx-yyy-zzz",
				}
				game.LastPickingUserId = nextUserId

				if _, ok := h.Rooms[message.RoomId]; ok {
					for _, client := range h.Rooms[message.RoomId].Clients {
						client.Message <- message
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
	//close(hub.Stall)
}

const WaitTime = 10 * time.Second

// Commands

const CommandUserJoined = "USER_JOINED"
const CommandBettingStarted = "BETTING_STARTED"
const CommandBettingEnded = "BETTING_ENDED"
const CommandDealCards = "DEAL_CARDS"
const CommandInitGame = "INIT_GAME"
const CommandPickingStarted = "PICKING_STARTED"

const CommandStart = "START"
const CommandPick = "PICK"
const CommandBet = "BET"

// Command Structs

type BetCommand struct {
	Round  int   `json:"round"`
	EndsAt int64 `json:"endsAt"`
}

type UserBet struct {
	UserId string `json:"userId"`
	Bet    int    `json:"bet"`
}

type CommandPickContent struct {
	CardId int `json:"cardId"`
}

type PickingStartedContent struct {
	UserId string `json:"userId"`
	EndsAt int64  `json:"endsAt"`
}
