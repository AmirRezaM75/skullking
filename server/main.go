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
	"skullking/pkg/validator"
	"skullking/repositories"
	"skullking/routes"
	"skullking/services"
)

func main() {
	loadEnvironments()

	client, cancel, disconnect := initDatabase()

	redis := initRedis()

	defer cancel()
	defer disconnect()

	db := client.Database(os.Getenv("MONGODB_DATABASE"))
	var userRepository = repositories.NewUserRepository(db)

	var tokenRepository = repositories.NewTokenRepository(redis)

	var userService = services.NewUserService(userRepository, tokenRepository)

	v := validator.NewValidator()

	r := router.NewRouter()
	r.Middleware(middlewares.CorsPolicy{})

	userHandler := handlers.NewUserHandler(userService, v)

	gameRepository := repositories.NewGameRepository(db)

	hub := models.NewHub(gameRepository)
	gameHandler := handlers.NewGameHandler(hub, userService)

	go hub.Run()

	routes.Route{
		Router:      r,
		UserService: userService,
		UserHandler: userHandler,
		GameHandler: gameHandler,
	}.Setup()

	fmt.Println("Listening on port 3000")

	log.Fatal(http.ListenAndServe(":3000", r))
}
