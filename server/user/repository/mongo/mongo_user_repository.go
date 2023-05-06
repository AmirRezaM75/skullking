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
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

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

func (ur mongoUserRepository) FindById(userId string) *domain.User {
	id, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return nil
	}

	filter := bson.M{"_id": id}

	var user domain.User

	err = ur.db.Collection(UsersTable).FindOne(context.Background(), filter).Decode(&user)

	if err != nil {
		return nil
	}

	return &user
}

func (ur mongoUserRepository) FindByEmail(email string) *domain.User {
	filter := bson.D{{"email", email}}

	var user domain.User

	err := ur.db.Collection(UsersTable).FindOne(context.Background(), filter).Decode(&user)

	if err != nil {
		return nil
	}

	return &user
}

func (ur mongoUserRepository) exists(filter bson.D) bool {
	count, err := ur.db.Collection(UsersTable).CountDocuments(
		context.Background(),
		filter,
		options.Count().SetLimit(1),
	)

	if err != nil {
		return false
	}

	return count != 0
}

func (ur mongoUserRepository) ExistsByUsername(username string) bool {
	return ur.exists(bson.D{{"username", username}})
}

func (ur mongoUserRepository) ExistsByEmail(email string) bool {
	return ur.exists(bson.D{{"email", email}})
}

func (ur mongoUserRepository) UpdateEmailVerifiedAtByUserId(userId string, datetime time.Time) bool {
	id, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return false
	}

	update := bson.D{{
		"$set",
		bson.D{{"email_verified_at", primitive.NewDateTimeFromTime(datetime)}},
	}}

	result, err := ur.db.Collection(UsersTable).UpdateByID(context.Background(), id, update)

	return result.ModifiedCount > 0
}
