package main

import (
	"context"
	"log"

	_ "github.com/lib/pq"
	"github.com/piriya-muaithaisong/testgolang_mongo/api"
	"github.com/piriya-muaithaisong/testgolang_mongo/db"
	"github.com/piriya-muaithaisong/testgolang_mongo/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config")
	}

	// if config.Environment == "development" {
	// 	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// }

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.DBSource))

	if err != nil {
		log.Fatal("cannot connect to the database")
	}

	// runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(client, config.DBName)

	// redisOpt := asynq.RedisClientOpt{
	// 	Addr: config.RedisAddress,
	// }

	// taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	runGinServer(config, store)
}

func runGinServer(config utils.Config, store *db.MongoStore) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ")
	}
}
