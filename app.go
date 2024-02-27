package skull_king

import (
	"context"
	"fmt"
	"github.com/amirrezam75/go-router"
	m "github.com/amirrezam75/go-router/middlewares"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"os"
	"skullking/handlers"
	"skullking/middlewares"
	"skullking/services"
	"time"
)

func initBroker() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("BROKER_REDIS_HOST"), os.Getenv("BROKER_REDIS_PORT")),
		Password: os.Getenv("BROKER_REDIS_PASSWORD"),
		DB:       0,
	})
}

func initDatabase() (*mongo.Client, context.CancelFunc, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var mongoURI = fmt.Sprintf("mongodb://%s:%d", os.Getenv("MONGODB_HOST"), 27017)

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

func setupRoutes(
	gameHandler *handlers.GameHandler,
	userService services.UserService,
) *router.Router {
	var rateLimiterLogger = func(identifier, url string) {
		services.LogService{}.Info(map[string]string{
			"message":    "Too many requests.",
			"identifier": identifier,
			"url":        url,
		})
	}

	var rateLimiterConfig = m.RateLimiterConfig{
		Duration: time.Minute,
		Limit:    10,
		Extractor: func(r *http.Request) string {
			var user = services.ContextService{}.GetUser(r.Context())

			if user == nil {
				return r.RemoteAddr
			}

			return user.Id
		},
	}

	var rateLimiterMiddleware = m.NewRateLimiterMiddleware(rateLimiterConfig, rateLimiterLogger)
	var authMiddleware = middlewares.Authenticate{UserService: userService}

	r := router.NewRouter()
	r.Middleware(middlewares.CorsPolicy{})

	r.Post("/games", gameHandler.Create).
		Middleware(authMiddleware).
		Middleware(rateLimiterMiddleware)

	r.Get("/games/join", gameHandler.Join)
	r.Get("/games/cards", gameHandler.Cards)

	return r
}
