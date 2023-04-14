package mongo

import (
	"context"
	"github.com/AmirRezaM75/skull-king/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	user.Id = result.InsertedID.(primitive.ObjectID)
	return &user, nil
}

func (ur mongoUserRepository) FindByUsername(username string) *domain.User {
	filter := bson.D{{"username", username}}

	var user domain.User

	err := ur.db.Collection(UsersTable).FindOne(context.Background(), filter).Decode(&user)

	if err != nil {
		return nil
	}

	return &user
}

func (ur mongoUserRepository) ExistsByUsername(username string) bool {
	count, err := ur.db.Collection(UsersTable).CountDocuments(
		context.Background(),
		bson.D{{"username", username}},
		options.Count().SetLimit(1),
	)

	if err != nil {
		return false
	}

	return count != 0
}

func (ur mongoUserRepository) ExistsByEmail(email string) bool {
	count, err := ur.db.Collection(UsersTable).CountDocuments(
		context.Background(),
		bson.D{{"email", email}},
		options.Count().SetLimit(1),
	)

	if err != nil {
		return false
	}

	return count != 0
}
