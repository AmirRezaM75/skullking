package handlers

import (
	"encoding/json"
	commonservices "github.com/amirrezam75/kenopsiacommon/services"
	"github.com/amirrezam75/kenopsiauser"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"net/http"
	"os"
	"skullking/constants"
	"skullking/models"
	"skullking/pkg/syncx"
	"skullking/responses"
	"skullking/services"
	"strconv"
	"time"
)

type GameHandler struct {
	hub            *models.Hub
	lobbyService   services.LobbyService
	userRepository kenopsiauser.UserRepository
}

func NewGameHandler(
	hub *models.Hub,
	lobbyService services.LobbyService,
	userRepository kenopsiauser.UserRepository,
) *GameHandler {
	return &GameHandler{
		hub:            hub,
		lobbyService:   lobbyService,
		userRepository: userRepository,
	}
}

func (gameHandler *GameHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := commonservices.ContextService{}.GetUser(r.Context())

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var payload struct {
		LobbyId string `json:"lobbyId"`
	}

	err := decoder(&payload, w, r)

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
		LobbyId:   lobby.Id,
	}

	rand.Seed(time.Now().UnixNano())

	var humansCount = len(lobby.Players)
	var botsCount = len(lobby.Bots)

	indexes := rand.Perm(humansCount + botsCount)

	var index = 0

	var players syncx.Map[string, *models.Player]

	for _, player := range lobby.Players {
		players.Store(player.Id, &models.Player{
			Id:          player.Id,
			Username:    player.Username,
			GameId:      game.Id,
			AvatarId:    player.AvatarId,
			Index:       indexes[index] + 1,
			IsConnected: false,
			IsClosed:    true,
			IsBot:       false,
		})
		index++
	}

	for _, bot := range lobby.Bots {
		var botId = strconv.Itoa(int(bot.Id))

		players.Store(botId, &models.Player{
			Id:          botId,
			Username:    bot.Username,
			GameId:      game.Id,
			AvatarId:    bot.AvatarId,
			Index:       indexes[index] + 1,
			IsConnected: true,
			IsClosed:    true,
			IsBot:       true,
		})
		index++
	}

	game.Players = players

	gameHandler.hub.Cleanup()
	gameHandler.hub.Games.Store(game.Id, game)

	message, err := responses.GameCreatedEvent(game.Id, game.LobbyId)

	if err != nil {
		services.LogService{}.Error(map[string]string{
			"message":     err.Error(),
			"description": "Can not marshal GameCreatedEvent",
			"method":      "GameHandler@Create",
		})
	}

	err = gameHandler.hub.PublisherService.Publish(message)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(
		responses.CreateGame{Id: game.Id},
	)

	if err != nil {
		services.LogService{}.Error(map[string]string{
			"message":     err.Error(),
			"description": "Can not marshal CreateGame response.",
			"method":      "GameHandler@Create",
		})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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
		services.LogService{}.Error(map[string]string{
			"message":     err.Error(),
			"method":      "GameHandler@Join",
			"description": "Upgrade TCP connection failed.",
		})
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

	userId, err := gameHandler.userRepository.AcquireUserId(ticketId)

	if err != nil {
		services.LogService{}.Error(map[string]string{
			"msg":  err.Error(),
			"desc": "could not acquire user id",
		})
		message := models.ServerMessage{
			Command: constants.CommandReportError,
			Content: responses.Error{
				Message:    "Unauthorized.",
				StatusCode: http.StatusUnauthorized,
			},
		}
		connection.WriteJSON(message)
		return
	}

	game, _ := gameHandler.hub.Games.Load(gameId)

	if player, exists := game.Players.Load(userId); exists {
		player.Kick()
		player.Message = make(chan *models.ServerMessage, 10)
		player.IsClosed = false
		player.Connection = connection
		player.IsConnected = true

		go player.Write()

		game.Initialize(gameHandler.hub, player.Id)

		if game.IsEveryoneConnected() == true &&
			game.State == constants.StatePending {
			game.Start(gameHandler.hub)
		}

		game.Joined(gameHandler.hub, player.Id)

		player.Read(gameHandler.hub)
	} else {
		message := models.ServerMessage{
			Command: constants.CommandReportError,
			Content: responses.Error{
				Message:    "You must join the game through lobby.",
				StatusCode: http.StatusForbidden,
			},
		}
		connection.WriteJSON(message)
	}
}
