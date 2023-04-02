package db

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserInt interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id string) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
}

type UserClient struct {
	client  *mongo.Client
	userCol *mongo.Collection
}

func NewUserClient(client *mongo.Client, dbName string) *UserClient {
	return &UserClient{
		client:  client,
		userCol: client.Database(dbName).Collection("user"),
	}
}

func (c *UserClient) CreateUser(ctx context.Context, user *User) error {
	_, err := c.userCol.InsertOne(ctx, user)
	if err != nil {
		log.Print(fmt.Errorf("could not add new user: %w", err))
		return err
	}
	return nil
}

func (c *UserClient) GetUserByID(ctx context.Context, id string) (User, error) {
	var User User
	objID, _ := primitive.ObjectIDFromHex(id)
	res := c.userCol.FindOne(ctx, bson.M{"_id": objID})

	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return User, nil
		}
		log.Print(fmt.Errorf("error when finding the User [%s]: %q", id, res.Err()))
		return User, res.Err()
	}

	if err := res.Decode(&User); err != nil {
		log.Print(fmt.Errorf("error decoding [%s]: %q", id, err))
		return User, err
	}

	return User, nil
}

func (c *UserClient) GetUserByUsername(ctx context.Context, username string) (User, error) {
	var User User
	res := c.userCol.FindOne(ctx, bson.M{"username": username})

	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return User, nil
		}
		log.Print(fmt.Errorf("error when finding the User [%s]: %q", username, res.Err()))
		return User, res.Err()
	}

	if err := res.Decode(&User); err != nil {
		log.Print(fmt.Errorf("error decoding [%s]: %q", username, err))
		return User, err
	}

	return User, nil
}
