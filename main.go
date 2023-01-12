package main

import (
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

func (s Server) start(w http.ResponseWriter, r *http.Request) {
	// This is not secure; but acceptable as starting point
	// @link: https://websockets.readthedocs.io/en/stable/topics/authentication.html
	token := r.URL.Query().Get("token")

	user := getUserByToken(token)

	_, ok := s.connections[user.id]

	if !ok {
		c, err := upgrader.Upgrade(w, r, nil)

		s.connections[user.id] = c

		fmt.Println("New Connection.")

		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		defer c.Close()
	} else {
		return
	}

	/*c, err := upgrader.Upgrade(w, r, nil)

	defer c.Close()

	s.connections[user.id] = c

	fmt.Println("New Connection.")

	if err != nil {
		log.Print("upgrade:", err)
		return
	}*/

	var deck Deck
	cards := generateCards()
	deck.cards = cards
	deck.shuffle()
	items := deck.deal(2, 3)

	fmt.Println("Number of connections: ", len(s.connections))
out:
	for {
		for userId, c := range s.connections {
			fmt.Println("UserId:", userId)
			err := c.WriteJSON(items[userId-1])
			if err != nil {
				log.Println("write:", err)
				delete(s.connections, userId)
				break out
			}
		}
		time.Sleep(time.Second * 5)
	}

	//if len(s.connections) == 2 {
	//
	//	for _, c := range s.connections {
	//		/*userCards, _ := json.Marshal(items[userId-1])
	//		stringJSON := string(userCards)
	//		err := c.WriteJSON(stringJSON)*/
	//		err = c.WriteMessage(websocket.TextMessage, []byte(time.Now().String()))
	//		if err != nil {
	//			log.Println("write:", err)
	//			break
	//		}
	//	}
	//}

	/*for {
		// The WebSocket protocol defines three types of control messages: close, ping and pong.
		// messageType is an int with value websocket.BinaryMessage (2) or websocket.TextMessage (1).
		messageType, message, err := c.ReadMessage() // or c.ReadJson
		if err != nil {
			log.Println("read:", err)
			break
		}

		log.Printf("message from client: %s", message)

		for _, c := range s.connections {
			err = c.WriteMessage(messageType, []byte(time.Now().String()))
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
		time.Sleep(time.Second * 2)
	}*/
}

func main() {

	server := Server{
		connections: make(map[int]*websocket.Conn),
	}
	http.HandleFunc("/start", server.start)
	fs := http.FileServer(http.Dir("client"))
	http.Handle("/", fs)
	log.Fatal(http.ListenAndServe(":3000", nil))
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
