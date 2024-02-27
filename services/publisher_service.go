package services

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
)

type PublisherService struct {
	broker *redis.Client
}

func NewPublisherService(broker *redis.Client) PublisherService {
	return PublisherService{broker: broker}
}

func (publisherService PublisherService) Publish(message string) error {
	if message == "" {
		return errors.New("message is empty")
	}

	err := publisherService.broker.Publish(context.Background(), "game-service", message).Err()

	if err != nil {
		LogService{}.Error(map[string]string{
			"method":  "PublisherService@publish",
			"message": err.Error(),
			"payload": message,
		})

		return err
	}

	return nil
}
