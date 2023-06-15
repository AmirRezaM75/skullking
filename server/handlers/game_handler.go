package handlers

import (
	"encoding/json"
	"github.com/AmirRezaM75/skull-king/constants"
	"github.com/AmirRezaM75/skull-king/contracts"
	"github.com/AmirRezaM75/skull-king/models"
	"github.com/AmirRezaM75/skull-king/pkg/support"
	"github.com/AmirRezaM75/skull-king/responses"
	"github.com/AmirRezaM75/skull-king/services"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
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

func (gameHandler *GameHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := services.ContextService{}.GetUser(r.Context())

	if user == nil {
		return
	}

	game := &models.Game{
		Id:        uuid.New().String(),
		State:     constants.StatePending,
		Players:   make(map[string]*models.Player, constants.MaxPlayers),
		Scores:    make(map[string]int, constants.MaxPlayers),
		CreatorId: user.Id.Hex(),
		CreatedAt: time.Now().Unix(),
	}

	gameHandler.hub.Cleanup()
	gameHandler.hub.Games[game.Id] = game

	response, err := json.Marshal(
		responses.CreateGame{Id: game.Id},
	)

	if err != nil {
		http.Error(w, "JSON marshal failed", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, _ = w.Write(response)
}

func (gameHandler *GameHandler) Cards(w http.ResponseWriter, _ *http.Request) {
	response := struct {
		Items []responses.Card `json:"items"`
	}{}

	for _, card := range models.GetCards() {
		response.Items = append(response.Items, responses.Card{
			Id:     int(card.Id),
			Number: card.Number,
			Type:   card.Type,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	return
}

func (gameHandler *GameHandler) Join(w http.ResponseWriter, r *http.Request) {
	gameId := r.URL.Query().Get("gameId")
	// TODO: The request can be intercepted
	// @link https://devcenter.heroku.com/articles/websocket-security
	token := r.URL.Query().Get("token")

	claims, err := support.ParseJWT(token)

	if err != nil {
		response := struct {
			Message string `json:"message"`
		}{Message: "Can not parse JWT."}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	user := gameHandler.userService.FindById(claims.ID)

	if user == nil {
		response := struct {
			Message string `json:"message"`
		}{Message: "User not found."}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	if user.EmailVerifiedAt == nil {
		response := struct {
			Message string `json:"message"`
		}{Message: "Email has not been verified."}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	if _, ok := gameHandler.hub.Games[gameId]; !ok {
		response := struct {
			Message string `json:"message"`
		}{Message: "Game not found."}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	game := gameHandler.hub.Games[gameId]

	if len(game.Players) == constants.MaxPlayers {
		response := struct {
			Message string `json:"message"`
		}{Message: "Game is already full."}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	connection, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		response := struct {
			Message string `json:"message"`
		}{Message: "Upgrade TCP connection failed."}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	var player *models.Player

	if _, exists := game.Players[user.Id.Hex()]; exists {
		player = game.Players[user.Id.Hex()]
		player.Connection = connection
		game.Initialize(gameHandler.hub, player.Id)
		return
	}

	player = &models.Player{
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
