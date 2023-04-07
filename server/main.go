package main

import (
	"context"
	"fmt"
	"github.com/AmirRezaM75/skull-king/user"
	"github.com/AmirRezaM75/skull-king/ws"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Load .env file

	loadEnv()

	// Setup Database

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	var mongoURI = fmt.Sprintf("mongodb://%s:%s", os.Getenv("MONGODB_HOST"), os.Getenv("MONGODB_PORT"))
	var credentials options.Credential = options.Credential{
		AuthSource: os.Getenv("MONGODB_AUTH_SOURCE"),
		Username:   os.Getenv("MONGODB_USERNAME"),
		Password:   os.Getenv("MONGODB_PASSWORD"),
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI).SetAuth(credentials))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		panic(err)
	}

	type Person struct {
		Name string
		Age  int
	}

	collection := client.Database(os.Getenv("MONGODB_DATABASE")).Collection("people")

	person := Person{"John Doe", 30}
	_, err = collection.InsertOne(context.Background(), person)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Inserted a document!")

	hub := ws.NewHub()

	wsHandler := ws.NewHandler(hub)

	go hub.Run()

	fs := http.FileServer(http.Dir("client"))

	http.Handle("/", fs)

	http.HandleFunc("/ws/join", wsHandler.Join)

	authHandler := user.NewHandler()

	http.HandleFunc("/register", authHandler.Register)

	fmt.Println("Listening on port 3000")

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func loadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
