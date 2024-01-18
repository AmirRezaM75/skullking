package handlers

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"math/rand"
	"net/http"
	"os"
	"skullking/constants"
	"skullking/models"
	"skullking/pkg/syncx"
	"skullking/responses"
	"skullking/services"
	"time"
)

type GameHandler struct {
	hub           *models.Hub
	lobbyService  services.LobbyService
	ticketService services.TicketService
}

func NewGameHandler(
	hub *models.Hub,
	lobbyService services.LobbyService,
	ticketService services.TicketService,
) *GameHandler {
	return &GameHandler{
		hub:           hub,
		lobbyService:  lobbyService,
		ticketService: ticketService,
	}
}

func (gameHandler *GameHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := services.ContextService{}.GetUser(r.Context())

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var payload struct {
		LobbyId string `json:"lobbyId"`
	}

	err := decoder(payload, w, r)

	if err != nil {
		return
	}

	var lobby = gameHandler.lobbyService.FindById(payload.LobbyId)

	if lobby == nil {
		errorResponse(w, "Lobby not found!", http.StatusNotFound)
		return
	}

	game := &models.Game{
		Id:        primitive.NewObjectID().Hex(),
		State:     constants.StatePending,
		CreatorId: user.Id,
		CreatedAt: time.Now().Unix(),
	}

	rand.Seed(time.Now().UnixNano())

	// Generate array of [0, players.length)
	indexes := rand.Perm(len(lobby.Players))

	var players syncx.Map[string, *models.Player]

	for i, player := range lobby.Players {
		players.Store(player.Id, &models.Player{
			Id:          player.Id,
			Username:    player.Username,
			GameId:      game.Id,
			AvatarId:    player.AvatarId,
			Message:     make(chan *models.ServerMessage, 10),
			Index:       indexes[i] + 1,
			IsConnected: false,
		})
	}

	game.Players = players

	gameHandler.hub.Cleanup()
	gameHandler.hub.Games.Store(game.Id, game)

	response, err := json.Marshal(
		responses.CreateGame{Id: game.Id},
	)

	if err != nil {
		http.Error(w, "JSON marshal failed", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(response)
}

func (gameHandler *GameHandler) Cards(w http.ResponseWriter, _ *http.Request) {
	response := struct {
		Items []responses.Card `json:"items"`
	}{}

	for _, card := range models.GetCards() {
		response.Items = append(response.Items, responses.Card{
			Id:     uint16(card.Id),
			Number: card.Number,
			Type:   card.Type,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	return
}

func (gameHandler *GameHandler) Join(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return r.Header.Get("origin") == os.Getenv("FRONTEND_URL")
		},
	}

	connection, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Upgrade TCP connection failed.", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	gameId := r.URL.Query().Get("gameId")

	if _, ok := gameHandler.hub.Games.Load(gameId); !ok {
		message := models.ServerMessage{
			Command: constants.CommandReportError,
			Content: responses.Error{
				Message:    "The game you are searching for could not be found. It's possible that the URL you entered is incorrect or the game has already concluded.",
				StatusCode: http.StatusNotFound,
			},
		}
		connection.WriteJSON(message)
		return
	}

	// @link https://devcenter.heroku.com/articles/websocket-security
	ticketId := r.URL.Query().Get("ticketId")

	if ticketId == "" {
		message := models.ServerMessage{
			Command: constants.CommandReportError,
			Content: responses.Error{
				Message:    "ticketId is required.",
				StatusCode: http.StatusUnprocessableEntity,
			},
		}
		connection.WriteJSON(message)
		return
	}

	userId := gameHandler.ticketService.AcquireUserId(ticketId)

	if userId == "" {
		message := models.ServerMessage{
			Command: constants.CommandReportError,
			Content: responses.Error{
				Message:    "Unauthorized.",
				StatusCode: http.StatusUnprocessableEntity,
			},
		}
		connection.WriteJSON(message)
		return
	}

	game, _ := gameHandler.hub.Games.Load(gameId)

	if player, exists := game.Players.Load(userId); exists {
		player.Connection = connection
		player.IsConnected = true

		go player.Write()

		game.Initialize(gameHandler.hub, player.Id)

		if game.IsEveryoneConnected() == true &&
			game.State == constants.StatePending {
			m := &models.ServerMessage{
				Command: constants.CommandStarted,
				GameId:  game.Id,
			}

			gameHandler.hub.Dispatch <- m

			game.Start(gameHandler.hub)
		}

		player.Read(gameHandler.hub)
	} else {
		message := models.ServerMessage{
			Command: constants.CommandReportError,
			Content: responses.Error{
				Message:    "You must join the game through lobby.",
				StatusCode: http.StatusBadRequest,
			},
		}
		connection.WriteJSON(message)
	}
}
