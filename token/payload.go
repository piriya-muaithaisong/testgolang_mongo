package token

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payload struct {
	ID        primitive.ObjectID `json:"id"`
	UserID    primitive.ObjectID `json:"user_id"`
	IssuedAt  time.Time          `json:"issued_at"`
	ExpiredAt time.Time          `json:"expired_at"`
}

func NewPayload(userID primitive.ObjectID, duration time.Duration) (*Payload, error) {
	tokenID := primitive.NewObjectID()

	payload := &Payload{
		ID:        tokenID,
		UserID:    userID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}
