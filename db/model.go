package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	RefreshToken string             `json:"refresh_token" bson:"refresh_token"`
	IsBlocked    bool               `json:"is_blocked" bson:"is_blocked"`
	ExpiresAt    time.Time          `json:"expires_at" bson:"expires_at"`
	CreateAt     time.Time          `json:"create_at" bson:"create_at"`
}

type User struct {
	ID                primitive.ObjectID `json:"id" bson:"_id"`
	Username          string             `json:"username" bson:"username"`
	HashedPassword    string             `json:"hashed_password" bson:"hashed_password"`
	FistName          string             `json:"first_name" bson:"first_name"`
	LastName          string             `json:"last_name" bson:"last_name"`
	Email             string             `json:"email" bson:"email"`
	CreateAt          time.Time          `json:"create_at" bson:"create_at"`
	PasswordChangedAt time.Time          `json:"password_changed_at" bson:"password_changed_at"`
	Org               string             `json:"org" bson:"org"`
	BillingAddress    string             `json:"billing_address" bson:"billing_address"`
	Role              string             `json:"role" bson:"role"`
}
