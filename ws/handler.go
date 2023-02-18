package ws

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

const MaxPlayers = 8

type Handler struct {
	hub *Hub
}

func NewHandler(h *Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

func (h *Handler) Create(w http.ResponseWriter, _ *http.Request) {

	gameId := uuid.New().String()

	h.hub.games[gameId] = &Game{
		id:      gameId,
		players: make(map[int]*Player, MaxPlayers),
	}

	response, err := json.Marshal(struct {
		Id string `json:"id"`
	}{
		Id: gameId,
	})

	if err != nil {
		http.Error(w, "JSON marshal failed", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, _ = w.Write(response)
}

func (h *Handler) Join(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		http.Error(w, "Upgrade TCP connection failed.", 500)
		return
	}

	gameId := r.URL.Query().Get("gameId")
	// TODO: Get token and find user from database + cache
	userId, _ := strconv.Atoi(r.URL.Query().Get("userId"))

	player := &Player{
		connection: c,
		message:    make(chan *ServerMessage, 10),
		id:         userId,
		gameId:     gameId,
		// Assign each player an avatar
	}

	content, _ := json.Marshal(struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
	}{
		Id:       userId,
		Username: fmt.Sprintf("Username #%d", userId),
	})

	h.hub.register <- player

	m := &ServerMessage{
		Command:     CommandJoined,
		ContentType: "json",
		Content:     string(content),
		gameId:      gameId,
		SenderId:    userId,
	}

	h.hub.dispatch <- m

	go player.writeMessage()
	player.readMessage(h.hub)
}
