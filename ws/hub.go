package ws

import (
	"time"
)

type Hub struct {
	games      map[string]*Game
	register   chan *Player
	unregister chan *Player
	dispatch   chan *ServerMessage
	remind     chan Reminder
}

type Reminder struct {
	gameId string
}

const StatePending = "PENDING"
const StateBidding = "BIDDING"
const StatePicking = "PICKING"
const StateCalculating = "CALCULATING"

func NewHub() *Hub {
	return &Hub{
		games:      make(map[string]*Game),
		register:   make(chan *Player),
		unregister: make(chan *Player),
		dispatch:   make(chan *ServerMessage),
		remind:     make(chan Reminder),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case reminder, ok := <-h.remind:
			game := h.games[reminder.gameId]
			if !ok {
				h.remind = nil
			} else {
				if game.state == StatePicking {
					endPicking(game, h)
				}
			}
		case player := <-h.register:
			if _, ok := h.games[player.gameId]; ok {
				h.games[player.gameId].players[player.id] = player
			}
		case _ = <-h.unregister:
			// Not going to remove player from game
			// because it should continue playing till the end
			// close(player.ServerMessage)
			// If every one left the game delete the game.

		case message := <-h.dispatch:
			// If there is no specific receiver broadcast it to all players
			if _, ok := h.games[message.gameId]; ok {
				var game = h.games[message.gameId]
				if message.receiverId == 0 {
					for _, player := range game.players {
						player.message <- message
					}
				} else {
					if _, ok = game.players[message.receiverId]; ok {
						game.players[message.receiverId].message <- message
					}
				}
			}
		}
	}
}

func (h *Hub) newReminder(gameId string) {
	time.Sleep(WaitTime)

	h.remind <- Reminder{
		gameId: gameId,
	}

	//close(hub.remind)
}

const WaitTime = 10 * time.Second

// Commands

const CommandJoined = "JOINED"
const CommandBiddingStarted = "BIDDING_STARTED"
const CommandBiddingEnded = "BIDDING_ENDED"
const CommandPickingStarted = "PICKING_STARTED"
const CommandPickingEnded = "PICKING_ENDED"
const CommandDeal = "DEAL"
const CommandInit = "INIT"
const CommandPicked = "PICKED"

// These commands coming from players

const CommandStart = "START"
const CommandPick = "PICK"
const CommandBid = "BID"

type BiddingStartedContent struct {
	Round  int   `json:"round"`
	EndsAt int64 `json:"endsAt"`
}

type UserBet struct {
	UserId string `json:"userId"`
	Bet    int    `json:"bet"`
}

type UserPickedCard struct {
	UserId string `json:"userId"`
	CardId int    `json:"cardId"`
}

type CommandPickContent struct {
	CardId int `json:"cardId"`
}

type PickingStartedContent struct {
	UserId int   `json:"userId"`
	EndsAt int64 `json:"endsAt"`
}
