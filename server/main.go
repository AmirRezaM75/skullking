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

	hub := models.NewHub(gameRepository)

	var broker = initBroker()

	var publisherService = services.NewPublisherService(broker)

	gameHandler := handlers.NewGameHandler(hub, lobbyService, ticketService, publisherService)

	go hub.Run()

	userService := services.NewUserService()

	var router = setupRoutes(gameHandler, userService)

	fmt.Println("Listening on port 3000")

	log.Fatal(http.ListenAndServe(":3000", router))
}
