package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id              primitive.ObjectID `bson:"_id,omitempty"`
	Username        string
	Email           string
	Password        string
	EmailVerifiedAt *primitive.DateTime `bson:"email_verified_at"`
	CreatedAt       primitive.DateTime  `bson:"created_at"`
}
