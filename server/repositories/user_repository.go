package repositories

import (
	"context"
	"github.com/AmirRezaM75/skull-king/contracts"
	"github.com/AmirRezaM75/skull-king/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const UsersTable = "users"

type userRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) contracts.UserRepository {
	return userRepository{
		db: db,
	}
}

func (ur userRepository) Create(user models.User) (*models.User, error) {
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	result, err := ur.db.Collection(UsersTable).InsertOne(context.Background(), user)

	if err != nil {
		return nil, err
	}

	user.Id = result.InsertedID.(primitive.ObjectID)
	return &user, nil
}

func (ur userRepository) FindByUsername(username string) *models.User {
	filter := bson.D{{"username", username}}

	var user models.User

	err := ur.db.Collection(UsersTable).FindOne(context.Background(), filter).Decode(&user)

	if err != nil {
		return nil
	}

	return &user
}

func (ur userRepository) FindById(userId string) *models.User {
	id, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return nil
	}

	filter := bson.M{"_id": id}

	var user models.User

	err = ur.db.Collection(UsersTable).FindOne(context.Background(), filter).Decode(&user)

	if err != nil {
		return nil
	}

	return &user
}

func (ur userRepository) FindByEmail(email string) *models.User {
	filter := bson.D{{"email", email}}

	var user models.User

	err := ur.db.Collection(UsersTable).FindOne(context.Background(), filter).Decode(&user)

	if err != nil {
		return nil
	}

	return &user
}

func (ur userRepository) exists(filter bson.D) bool {
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

func (ur userRepository) ExistsByUsername(username string) bool {
	return ur.exists(bson.D{{"username", username}})
}

func (ur userRepository) ExistsByEmail(email string) bool {
	return ur.exists(bson.D{{"email", email}})
}

func (ur userRepository) UpdateEmailVerifiedAtByUserId(userId string, datetime time.Time) bool {
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

func (ur userRepository) UpdatePasswordByEmail(email, password string) bool {
	filter := bson.M{"email": email}

	update := bson.M{"$set": bson.M{"password": password}}

	result, err := ur.db.Collection(UsersTable).UpdateOne(context.Background(), filter, update)

	if err != nil {
		return false
	}

	return result.ModifiedCount > 0
}
