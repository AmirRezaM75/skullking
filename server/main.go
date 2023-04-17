package main

import (
	"bytes"
	"fmt"
	"github.com/AmirRezaM75/skull-king/app"
	"github.com/AmirRezaM75/skull-king/pkg/support"
	"github.com/AmirRezaM75/skull-king/pkg/validator"
	_userHandler "github.com/AmirRezaM75/skull-king/user/delivery/http"
	_userRepository "github.com/AmirRezaM75/skull-king/user/repository/mongo"
	_userService "github.com/AmirRezaM75/skull-king/user/service"
	"github.com/AmirRezaM75/skull-king/ws"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	application := app.App{}

	application.LoadEnvironments()

	var body bytes.Buffer
	t, err := template.ParseFiles("index.html")

	if err != nil {
		fmt.Println("Can't parse HTML file.")
	}

	t.Execute(&body, nil)

	m := support.Mail{
		To:      []string{"amir@gmail.com"},
		Subject: "Register",
		Body:    body.String(),
	}

	m.Send()

	client, cancel, disconnect := application.InitDatabase()

	defer cancel()
	defer disconnect()

	var userRepository = _userRepository.NewMongoUserRepository(
		client.Database(os.Getenv("MONGODB_DATABASE")),
	)

	var userService = _userService.NewUserService(userRepository)

	v := validator.NewValidator()

	_userHandler.NewUserHandler(userService, v)

	hub := ws.NewHub()

	wsHandler := ws.NewHandler(hub)

	go hub.Run()

	http.HandleFunc("/ws/join", wsHandler.Join)

	fmt.Println("Listening on port 3000")

	log.Fatal(http.ListenAndServe(":3000", nil))
}
