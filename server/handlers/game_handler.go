package handlers

import (
	"encoding/json"
	"github.com/AmirRezaM75/skull-king/constants"
	"github.com/AmirRezaM75/skull-king/contracts"
	"github.com/AmirRezaM75/skull-king/models"
	"github.com/AmirRezaM75/skull-king/pkg/support"
	"github.com/AmirRezaM75/skull-king/responses"
	"github.com/AmirRezaM75/skull-king/services"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"os"
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
		Id:        primitive.NewObjectID().Hex(),
		State:     constants.StatePending,
		CreatorId: user.Id.Hex(),
		CreatedAt: time.Now().Unix(),
	}

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
	w.WriteHeader(200)
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

	// TODO: The request can be intercepted
	// @link https://devcenter.heroku.com/articles/websocket-security
	token := r.URL.Query().Get("token")

	claims, err := support.ParseJWT(token)

	if err != nil {
		message := models.ServerMessage{
			Command: constants.CommandReportError,
			Content: responses.Error{
				Message:    "Can not parse JWT.",
				StatusCode: http.StatusUnauthorized,
			},
		}
		connection.WriteJSON(message)
		return
	}

	// TODO: Why not using middleware
	user := gameHandler.userService.FindById(claims.ID)

	if user == nil {
		message := models.ServerMessage{
			Command: constants.CommandReportError,
			Content: responses.Error{
				Message:    "User not found.",
				StatusCode: http.StatusNotFound,
			},
		}
		connection.WriteJSON(message)
		return
	}

	if user.EmailVerifiedAt == nil {
		message := models.ServerMessage{
			Command: constants.CommandReportError,
			Content: responses.Error{
				Message:    "Email has not been verified.",
				StatusCode: http.StatusUnauthorized,
			},
		}
		connection.WriteJSON(message)
		return
	}

	game, _ := gameHandler.hub.Games.Load(gameId)

	if _, exists := game.Players.Load(user.Id.Hex()); !exists &&
		game.Players.Len() == constants.MaxPlayers {
		message := models.ServerMessage{
			Command: constants.CommandReportError,
			Content: responses.Error{
				Message:    "Game is already full.",
				StatusCode: http.StatusBadRequest,
			},
		}
		connection.WriteJSON(message)
		return
	}

	var player *models.Player

	if _, exists := game.Players.Load(user.Id.Hex()); exists {
		player, _ = game.Players.Load(user.Id.Hex())
		player.Connection = connection
		player.Message = make(chan *models.ServerMessage, 10)
		player.IsConnected = true
		game.Initialize(gameHandler.hub, player.Id)
	} else {
		if game.State != constants.StatePending {
			message := models.ServerMessage{
				Command: constants.CommandReportError,
				Content: responses.Error{
					Message:    "You can not join the game that has already started.",
					StatusCode: http.StatusBadRequest,
				},
			}
			connection.WriteJSON(message)
			return
		}

		player = &models.Player{
			Id:          user.Id.Hex(),
			Username:    user.Username,
			GameId:      gameId,
			Avatar:      game.GetAvailableAvatar(),
			Connection:  connection,
			Message:     make(chan *models.ServerMessage, 10),
			Index:       int(time.Now().UnixMilli()),
			IsConnected: true,
		}

		gameHandler.hub.Subscribe(player)

		game.Initialize(gameHandler.hub, player.Id)

		game.Join(gameHandler.hub, player)
	}

	go player.Write()

	player.Read(gameHandler.hub)
}
