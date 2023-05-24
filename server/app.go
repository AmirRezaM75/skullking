package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

func initDatabase() (*mongo.Client, context.CancelFunc, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var mongoURI = fmt.Sprintf("mongodb://%s:%s", os.Getenv("MONGODB_HOST"), os.Getenv("MONGODB_PORT"))

	var credentials = options.Credential{
		AuthSource: os.Getenv("MONGODB_AUTH_SOURCE"),
		Username:   os.Getenv("MONGODB_USERNAME"),
		Password:   os.Getenv("MONGODB_PASSWORD"),
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI).SetAuth(credentials))

	var disconnect = func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}

	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		panic(err)
	}

	return client, cancel, disconnect
}

func loadEnvironments() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func initRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
}