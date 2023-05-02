package configs

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

func ConnectDB() {

	fmt.Println("Loading mongo uri" + Config.MongoUri)
	client, err := mongo.NewClient(options.Client().ApplyURI(Config.MongoUri))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to mongodb")
	MI = MongoInstance{
		Client: client,
	}

}

type MongoInstance struct {
	Client *mongo.Client
}

var MI MongoInstance

func GetCollection(collectionName string) *mongo.Collection {
	collection := MI.Client.Database(Config.MongoDatabase).Collection(collectionName)
	return collection
}
