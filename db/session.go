package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

type SessionInt interface {
	CreateSession(ctx context.Context, session *Session) error
}

type SessionClient struct {
	client     *mongo.Client
	sessionCol *mongo.Collection
}

func NewSessionClient(client *mongo.Client, dbName string) *SessionClient {
	return &SessionClient{
		client:     client,
		sessionCol: client.Database(dbName).Collection("session"),
	}
}

func (c *SessionClient) CreateSession(ctx context.Context, session *Session) error {
	_, err := c.sessionCol.InsertOne(ctx, session)
	if err != nil {
		log.Print(fmt.Errorf("could not add new session: %w", err))
		return err
	}
	return nil
}

// func (c *SessionClient) GetSession(ctx context.Context,  string) error {
// }
