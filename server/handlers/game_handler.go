package handlers

import (
	"encoding/json"
	"github.com/AmirRezaM75/skull-king/constants"
	"github.com/AmirRezaM75/skull-king/contracts"
	"github.com/AmirRezaM75/skull-king/models"
	"github.com/AmirRezaM75/skull-king/pkg/support"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
)

type GameHandler struct {
	hub         *models.Hub
	userService contracts.UserService
}

func NewGameHandler(hub *models.Hub, userService contracts.UserService) *GameHandler {
	return &GameHandler{
		hub:         hub,
		userService: userService,
	}
}

func (gameHandler *GameHandler) Create(w http.ResponseWriter, _ *http.Request) {

	gameId := uuid.New().String()

	gameHandler.hub.Games[gameId] = &models.Game{
		Id:      gameId,
		State:   constants.StatePending,
		Players: make(map[string]*models.Player, constants.MaxPlayers),
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

func (gameHandler *GameHandler) Join(w http.ResponseWriter, r *http.Request) {
	gameId := r.URL.Query().Get("gameId")
	// TODO: The request can be intercepted
	// @link https://devcenter.heroku.com/articles/websocket-security
	token := r.URL.Query().Get("token")

	claims, err := support.ParseJWT(token)

	if err != nil {
		r := struct {
			Message string `json:"message"`
		}{Message: "Can not parse JWT."}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(r)
		return
	}

	user := gameHandler.userService.FindById(claims.ID)

	if user == nil {
		r := struct {
			Message string `json:"message"`
		}{Message: "User not found."}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(r)
		return
	}

	if user.EmailVerifiedAt == nil {
		r := struct {
			Message string `json:"message"`
		}{Message: "Email has not been verified."}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(r)
		return
	}

	if _, ok := gameHandler.hub.Games[gameId]; !ok {
		r := struct {
			Message string `json:"message"`
		}{Message: "Game not found."}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(r)
		return
	}

	game := gameHandler.hub.Games[gameId]

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	connection, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		r := struct {
			Message string `json:"message"`
		}{Message: "Upgrade TCP connection failed."}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(r)
		return
	}

	player := &models.Player{
		Id:         user.Id.Hex(),
		Username:   user.Username,
		GameId:     gameId,
		Avatar:     game.GetAvailableAvatar(),
		Connection: connection,
		Message:    make(chan *models.ServerMessage, 10),
	}

	gameHandler.hub.Subscribe(player)

	game.Initialize(gameHandler.hub, player.Id)

	game.Join(gameHandler.hub, player)

	go player.Write()

	player.Read(gameHandler.hub)
}
