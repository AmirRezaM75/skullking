package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"time"
)

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

/*
	count: number of players
	size: number of cards associated to each player
*/
//TODO: no need to count
func (d Deck) deal(count, size int) [][]Card {
	var output [][]Card

	index := 0

	for i := 0; i < count; i++ {
		output = append(output, d.cards[index:size+index])
		index = size + index
	}

	return output
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Server struct {
	connections map[int]*websocket.Conn
}

type User struct {
	id    int
	token string
}

type BetCommand struct {
	Method string
	EndsAt int64
}

func (s Server) start(w http.ResponseWriter, r *http.Request) {
	//index := 1
	// This is not secure; but acceptable as starting point
	// @link: https://websockets.readthedocs.io/en/stable/topics/authentication.html
	token := r.URL.Query().Get("token")

	user := getUserByToken(token)

	c, err := upgrader.Upgrade(w, r, nil)

	s.connections[user.id] = c

	fmt.Println("New Connection.")

	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	bets := make(map[int]int)

	defer c.Close()
	var betting = make(chan bool)

	for {
		reader := make(chan string)
		go s.reader(c, reader)
		token = r.URL.Query().Get("token")
		user = getUserByToken(token)

		select {
		case command := <-reader:
			fmt.Println("Received message from userId:", user.id, command)
			if command == "start" {
				log.Println("Started")
				endsAt := time.Now().Add(time.Second * 10).Unix()
				command, _ := json.Marshal(BetCommand{
					Method: "bet",
					EndsAt: endsAt,
				})
				for _, c := range s.connections {
					stringJSON := string(command)
					err := c.WriteJSON(stringJSON)
					if err != nil {
						log.Println("write:", err)
						break
					}
				}

				var deck Deck
				cards := generateCards()
				deck.cards = cards
				deck.shuffle()
				items := deck.deal(len(s.connections), 3)

				for userId, c := range s.connections {
					userCards, _ := json.Marshal(items[userId-1])
					stringJSON := string(userCards)
					err := c.WriteJSON(stringJSON)
					if err != nil {
						log.Println("write:", err)
						break
					}
				}
			}
			if command == "start" {
				go wait(betting)
			}
			if command == "bet" {
				bets[1] = 12
			}
		case expired := <-betting:
			if expired {
				fmt.Println("connections", len(s.connections))
				for userId, c := range s.connections {
					_, exists := bets[userId]
					if !exists {
						bets[userId] = 2
					}
					err := c.WriteJSON(bets)
					if err != nil {
						log.Println("write:", err)
						break
					}
				}
			}
		}
	}
}

func (s Server) reader(c *websocket.Conn, channel chan string) {
	_, message, err := c.ReadMessage()

	if err != nil {
		log.Println("ReadMessage:", err)
		channel <- ""
	}

	channel <- string(message)

	close(channel)
}

func wait(channel chan bool) {
	time.Sleep(time.Second * 5)
	channel <- true
	close(channel)
}

func main() {

	server := Server{
		connections: make(map[int]*websocket.Conn),
	}
	http.HandleFunc("/start", server.start)
	fs := http.FileServer(http.Dir("client"))
	http.Handle("/", fs)
	http.ListenAndServe(":3000", nil)
}

// TODO: Get user id from database by token
func getUserByToken(token string) User {
	userId := 1

	if token == "POORLY_SECURE_TOKEN" {
		userId = 2
	}

	return User{
		id:    userId,
		token: token,
	}
}
