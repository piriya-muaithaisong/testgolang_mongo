package token

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrExpiredToken = errors.New("token has expired")
var ErrInvalidToken = errors.New("token is invalid")

type Maker interface {
	CreateToken(userID primitive.ObjectID, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
