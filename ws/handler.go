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
	// TODO: Get token and find user from database + cache
	userId := r.URL.Query().Get("userId")

	client := &Client{
		Connection: c,
		Message:    make(chan *Message, 10),
		Id:         userId,
		RoomId:     roomId,
	}

	user := ClientRes{
		Id:       userId,
		Username: "Username #" + userId,
	}

	userBytes, _ := json.Marshal(user)

	m := &Message{
		Command:     CommandUserJoined,
		ContentType: "json",
		Content:     string(userBytes),
		RoomId:      roomId,
		SenderId:    client.Id,
	}

	h.hub.Register <- client
	h.hub.Broadcast <- m

	//TODO: Get game state from database
	var clients []ClientRes
	game := h.hub.Rooms[client.RoomId]
	for _, clientRes := range game.Clients {
		clients = append(clients, ClientRes{
			Id:       clientRes.Id,
			Username: "Username #" + clientRes.Id,
			Bid:      clientRes.Bid,
		})
	}
	gameResponse := GameRes{
		Round:  game.Round,
		Status: game.Status,
		Users:  clients,
	}
	gameResponseBytes, _ := json.Marshal(gameResponse)
	m = &Message{
		Command:     CommandInitGame,
		ContentType: "json",
		Content:     string(gameResponseBytes),
		RoomId:      roomId,
		ReceiverId:  client.Id,
	}

	h.hub.Broadcast <- m

	go client.writeMessage()
	client.readMessage(h.hub)
}

type RoomRes struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) GetRooms(w http.ResponseWriter, _ *http.Request) {
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

type GameRes struct {
	Round  int         `json:"round"`
	Status string      `json:"status"`
	Users  []ClientRes `json:"users"`
}

type ClientRes struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Bid      int    `json:"bid"`
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
