package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("DB_URI")))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	defer cancel()

	if err != nil {
		log.Fatal(err)
	}

	DB = client
	return client
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("mydb").Collection(collectionName)
}
