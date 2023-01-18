package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type Handler struct {
	hub *Hub
}

func NewHandler(h *Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

type CreateRoomRequest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var p struct {
		Id   string
		Name string
	}

	err := json.NewDecoder(r.Body).Decode(&p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	request := CreateRoomRequest{
		Id:   p.Id,
		Name: p.Name,
	}

	h.hub.Rooms[request.Id] = &Room{
		Id:      request.Id,
		Name:    request.Name,
		Clients: make(map[string]*Client),
	}

	bytes, err := json.Marshal(request)

	if err != nil {
		http.Error(w, "JSON marshal failed", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, _ = w.Write(bytes)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		http.Error(w, "Upgrade TCP connection failed.", 500)
		return
	}

	roomId := r.URL.Query().Get("roomId")
	userId := r.URL.Query().Get("userId")

	client := &Client{
		Connection: c,
		Message:    make(chan *Message, 10),
		Id:         userId,
		RoomId:     roomId,
	}

	m := &Message{
		Content: fmt.Sprintf("user %s joined the room.", userId),
		RoomId:  roomId,
	}

	h.hub.Register <- client
	h.hub.Broadcast <- m

	go client.writeMessage()
	client.readMessage(h.hub)
}

type RoomRes struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) GetRooms(w http.ResponseWriter, r *http.Request) {
	rooms := make([]RoomRes, 0)

	for _, r := range h.hub.Rooms {
		rooms = append(rooms, RoomRes{
			Id:   r.Id,
			Name: r.Name,
		})
	}

	bytes, err := json.Marshal(rooms)

	if err != nil {
		http.Error(w, "JSON marshal failed", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, _ = w.Write(bytes)
}

type ClientRes struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

func (h *Handler) GetClients(w http.ResponseWriter, r *http.Request) {
	var clients []ClientRes
	roomId := r.URL.Query().Get("roomId")

	if _, ok := h.hub.Rooms[roomId]; !ok {
		http.NotFound(w, r)
	}

	for _, c := range h.hub.Rooms[roomId].Clients {
		clients = append(clients, ClientRes{
			Id: c.Id,
		})
	}
	fmt.Println(clients)
	bytes, err := json.Marshal(clients)

	if err != nil {
		http.Error(w, "JSON marshal failed", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, _ = w.Write(bytes)
}
