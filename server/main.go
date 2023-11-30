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

	userService := services.NewUserService()

	gameRepository := repositories.NewGameRepository(db)

	hub := models.NewHub(gameRepository)
	gameHandler := handlers.NewGameHandler(hub, userService)

	go hub.Run()

	r.Post("/games", gameHandler.Create).
		Middleware(middlewares.Authenticate{UserService: userService})
	r.Get("/games/join", gameHandler.Join)
	r.Get("/games/cards", gameHandler.Cards)

	fmt.Println("Listening on port 3000")

	log.Fatal(http.ListenAndServe(":3000", r))
}
