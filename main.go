package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"skullking/handlers"
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

	var ticketService = services.NewTicketService(os.Getenv("KENOPSIA_USER_BASE_URL"))

	var lobbyService = services.NewLobbyService(os.Getenv("KENOPSIA_LOBBY_BASE_URL"), os.Getenv("KENOPSIA_TOKEN"))

	gameRepository := repositories.NewGameRepository(db)

	botRepository := repositories.NewBotRepository(os.Getenv("SKULLKING_AI_BASE_URL"))

	var broker = initBroker()

	var publisherService = services.NewPublisherService(broker)

	var logService = services.LogService{}

	hub := models.NewHub(gameRepository, botRepository, publisherService, logService)

	gameHandler := handlers.NewGameHandler(hub, lobbyService, ticketService)

	go hub.Run()

	userService := services.NewUserService()

	var router = setupRoutes(gameHandler, userService)

	fmt.Println("Listening on port 3000")

	log.Fatal(http.ListenAndServe(":3000", router))
}
