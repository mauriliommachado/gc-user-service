package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User model
type User struct {
	ID       primitive.ObjectID `json:"id"`
	Username string             `json:"username"`
	Password string             `json:"password"`
}
