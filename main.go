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
	color  string
	number int
}

func generateCards() []Card {
	colors := [4]string{"Black", "Red", "Green", "Yellow"}

	var cards []Card

	for _, color := range colors {
		for i := 1; i <= 14; i++ {
			card := Card{
				color:  color,
				number: i,
			}
			cards = append(cards, card)
		}
	}

	return cards
}

var upgrader = websocket.Upgrader{}

func start(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		// The WebSocket protocol defines three types of control messages: close, ping and pong.

		// messageType is an int with value websocket.BinaryMessage (2) or websocket.TextMessage (1).
		messageType, message, err := c.ReadMessage() // or c.ReadJson
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		err = c.WriteMessage(messageType, []byte(time.Now().String()))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	var deck Deck
	cards := generateCards()
	deck.cards = cards
	deck.shuffle()
	fmt.Println(deck.cards)
	output := deck.deal(4, 3)
	fmt.Println(output)
	return

	http.HandleFunc("/start", start)
	fs := http.FileServer(http.Dir("client"))
	http.Handle("/", fs)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
