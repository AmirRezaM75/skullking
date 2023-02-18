package main

import (
	"fmt"
	"github.com/AmirRezaM75/skull-king/ws"
	"log"
	"net/http"
)

func main() {
	hub := ws.NewHub()

	wsHandler := ws.NewHandler(hub)

	go hub.Run()

	fs := http.FileServer(http.Dir("client"))

	http.Handle("/", fs)

	http.HandleFunc("/ws/join", wsHandler.Join)

	fmt.Println("Listening on port 3000")

	log.Fatal(http.ListenAndServe(":3000", nil))
}
