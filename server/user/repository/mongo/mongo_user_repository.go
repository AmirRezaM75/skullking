package mongo

import (
	"context"
	"github.com/AmirRezaM75/skull-king/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const UsersTable = "users"

type mongoUserRepository struct {
	db *mongo.Database
}

func NewMongoUserRepository(db *mongo.Database) domain.UserRepository {
	return mongoUserRepository{
		db: db,
	}
}

func (ur mongoUserRepository) Create(user domain.User) (*domain.User, error) {

	user.CreatedAt = int64(primitive.NewDateTimeFromTime(time.Now()))

	result, err := ur.db.Collection(UsersTable).InsertOne(context.Background(), user)

	if err != nil {
		return nil, err
	}

	user.Id = result.InsertedID.(primitive.ObjectID).Hex()
	return &user, nil
}
