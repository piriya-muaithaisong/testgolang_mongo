package db

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoStore struct {
	Client        *mongo.Client
	SessionClient *SessionClient
	UserClient    *UserClient
}

func NewStore(client *mongo.Client, dbName string) *MongoStore {

	return &MongoStore{
		Client:        client,
		SessionClient: NewSessionClient(client, dbName),
		UserClient:    NewUserClient(client, dbName),
	}
}
