package main

import (
	"fmt"
	"github.com/amirrezam75/go-router"
	"log"
	"net/http"
	"os"
	"skullking/handlers"
	"skullking/middlewares"
	"skullking/models"
	"skullking/repositories"
	"skullking/services"
)

func main() {
	loadEnvironments()

	client, cancel, disconnect := initDatabase()

	defer cancel()
	defer disconnect()

	db := client.Database(os.Getenv("MONGODB_DATABASE"))

	r := router.NewRouter()
	r.Middleware(middlewares.CorsPolicy{})

	var ticketService = services.NewTicketService(os.Getenv("KENOPSIA_USER_BASE_URL"))

	var lobbyService = services.NewLobbyService(os.Getenv("KENOPSIA_LOBBY_BASE_URL"), os.Getenv("KENOPSIA_TOKEN"))

	userService := services.NewUserService()

	gameRepository := repositories.NewGameRepository(db)

	hub := models.NewHub(gameRepository)

	var broker = initBroker()

	var publisherService = services.NewPublisherService(broker)

	gameHandler := handlers.NewGameHandler(hub, lobbyService, ticketService, publisherService)

	go hub.Run()

	r.Post("/games", gameHandler.Create).
		Middleware(middlewares.Authenticate{UserService: userService})
	r.Get("/games/join", gameHandler.Join)
	r.Get("/games/cards", gameHandler.Cards)

	fmt.Println("Listening on port 3000")

	log.Fatal(http.ListenAndServe(":3000", r))
}
