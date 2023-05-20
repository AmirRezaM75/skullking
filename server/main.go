package main

import (
	"fmt"
	ws "github.com/AmirRezaM75/skull-king/game/hub"
	"github.com/AmirRezaM75/skull-king/handlers"
	"github.com/AmirRezaM75/skull-king/middlewares"
	"github.com/AmirRezaM75/skull-king/pkg/router"
	"github.com/AmirRezaM75/skull-king/pkg/validator"
	"github.com/AmirRezaM75/skull-king/repositories"
	"github.com/AmirRezaM75/skull-king/services"
	"log"
	"net/http"
	"os"
)

func main() {
	loadEnvironments()

	client, cancel, disconnect := initDatabase()

	redis := initRedis()

	defer cancel()
	defer disconnect()

	var userRepository = repositories.NewUserRepository(
		client.Database(os.Getenv("MONGODB_DATABASE")),
	)

	var tokenRepository = repositories.NewTokenRepository(redis)

	var userService = services.NewUserService(userRepository, tokenRepository)

	v := validator.NewValidator()

	r := router.NewRouter()
	r.Middleware(middlewares.CorsPolicy{})

	handlers.NewUserHandler(userService, v, r)

	hub := ws.NewHub()

	wsHandler := ws.NewHandler(hub)

	go hub.Run()

	http.HandleFunc("/ws/join", wsHandler.Join)

	fmt.Println("Listening on port 3000")

	log.Fatal(http.ListenAndServe(":3000", r))
}
