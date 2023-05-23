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
	// TODO: The request can be intercepted
	// @link https://devcenter.heroku.com/articles/websocket-security
	token := r.URL.Query().Get("token")

	claims, err := support.ParseJWT(token)

	if err != nil {
		// TODO:
		return
	}

	user := gameHandler.userService.FindById(claims.ID)

	player := &models.Player{
		Connection: c,
		Message:    make(chan *models.ServerMessage, 10),
		Id:         user.Id.Hex(),
		GameId:     gameId,
		// TODO: Assign each player an avatar
	}

	content, _ := json.Marshal(struct {
		Id       string `json:"id"`
		Username string `json:"username"`
	}{
		Id:       user.Id.Hex(),
		Username: user.Username,
	})

	gameHandler.hub.Register <- player

	m := &models.ServerMessage{
		Command:     constants.CommandJoined,
		ContentType: "json",
		Content:     string(content),
		GameId:      gameId,
		SenderId:    user.Id.Hex(),
	}

	gameHandler.hub.Dispatch <- m

	go player.Write()
	player.Read(gameHandler.hub)
}
